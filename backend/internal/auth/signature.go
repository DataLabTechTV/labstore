package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/DataLabTechTV/labstore/backend/internal/security"
	"github.com/DataLabTechTV/labstore/backend/pkg/iam"
)

const UnsignedPayload = "UNSIGNED-PAYLOAD"

type SigV4Result struct {
	AccessKey   string
	SecretKey   string
	Signature   string
	Timestamp   string
	Scope       string
	IsStreaming bool
}

func VerifySigV4(r *http.Request) (*SigV4Result, error) {
	// !FIXME: Could we refactor this into a few more functions?

	auth := r.Header.Get("Authorization")
	slog.Debug("Processing SigV4 request", "Authorization", security.TruncParamHeader(auth, "Signature"))

	payloadHash := r.Header.Get("X-Amz-Content-SHA256")
	slog.Debug("Received payload hash", "X-Amz-Content-SHA256", security.Trunc(payloadHash))

	// Remove prefix

	if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256") {
		return nil, errors.New("header Authorization must start with AWS4-HMAC-SHA256")
	}

	auth, ok := strings.CutPrefix(auth, "AWS4-HMAC-SHA256 ")
	if !ok {
		return nil, errors.New("could not remove prefix AWS4-HMAC-SHA256")
	}

	// Parse credentials, signed headers, and signature

	parts := strings.Split(auth, ",")

	var credential string
	var signedHeaders []string
	var signature string

	for _, p := range parts {
		p = strings.TrimSpace(p)

		if after, ok := strings.CutPrefix(p, "Credential="); ok {
			credential = after
		}

		if after, ok := strings.CutPrefix(p, "SignedHeaders="); ok {
			signedHeaders = strings.Split(after, ";")
		}

		if after, ok := strings.CutPrefix(p, "Signature="); ok {
			signature = after
		}
	}

	if credential == "" {
		return nil, errors.New("header Credential is empty")
	}

	if len(signedHeaders) == 0 {
		return nil, errors.New("header SignedHeaders is empty")
	}

	if signature == "" {
		return nil, errors.New("header Signature is empty")
	}

	slog.Debug(
		"Extracted authorization header parts",
		"Credential", credential,
		"SignedHeaders", strings.Join(signedHeaders, ";"),
		"Signature", security.Trunc(signature),
	)

	// Extract access key and scope

	credentialParts := strings.Split(credential, "/")

	accessKey := credentialParts[0]
	slog.Debug("Extracted access key from credential", "accessKey", accessKey)

	secretKey, ok := iam.Users[accessKey]
	if !ok {
		return nil, fmt.Errorf("no secret key found for access key %s", accessKey)
	}

	scope := strings.Join(credentialParts[1:], "/")
	slog.Debug("Extracted scope from credential", "scope", scope)

	// Compute signature

	canonicalRequest, err := buildCanonicalRequest(r, signedHeaders, payloadHash)
	if err != nil {
		return nil, errors.New("could not build canonical request")
	}
	slog.Debug("Built canonical request", "canonicalRequest", security.TruncLastLine(canonicalRequest))

	timestamp := r.Header.Get("X-Amz-Date")
	slog.Debug("Received timestamp", "X-Amz-Date", timestamp)

	stringToSign := buildStringToSign(timestamp, scope, canonicalRequest)
	slog.Debug("Built string to sign", "stringToSign", security.TruncLastLine(stringToSign))

	recomputedSignature, err := computeSignature(secretKey, scope, stringToSign)
	if err != nil {
		return nil, errors.New("could not compute signature")
	}

	byteSignature, err := hex.DecodeString(signature)
	if err != nil {
		return nil, errors.New("could not decode original signature")
	}

	byteRecomputedSignature, err := hex.DecodeString((recomputedSignature))
	if err != nil {
		return nil, errors.New("could not decode recomputed signature")
	}

	slog.Debug(
		"Comparing signatures",
		"original", security.Trunc(signature),
		"recomputed", security.Trunc(recomputedSignature),
	)

	if hmac.Equal(byteSignature, byteRecomputedSignature) {
		isStreaming := payloadHash == StreamingPayload

		res := &SigV4Result{
			AccessKey:   accessKey,
			SecretKey:   secretKey,
			Signature:   recomputedSignature,
			Timestamp:   timestamp,
			Scope:       scope,
			IsStreaming: isStreaming,
		}

		return res, nil
	}

	slog.Error("Original and recomputed signatures differ")
	return nil, errors.New("signatures do not match")
}

func buildCanonicalRequest(
	r *http.Request,
	signedHeaders []string,
	payloadHash string,
) (string, error) {
	var canonicalRequest strings.Builder

	canonicalRequest.WriteString(r.Method)
	canonicalRequest.WriteString("\n")

	uri := buildCanonicalURI(r.URL.Path)
	slog.Debug("Built canonical URI", "uri", uri)
	canonicalRequest.WriteString(uri)
	canonicalRequest.WriteString("\n")

	queryString := buildQueryString(r.URL.RawQuery)
	slog.Debug("Built canonical query string", "queryString", queryString)
	canonicalRequest.WriteString(queryString)
	canonicalRequest.WriteString("\n")

	for _, signedHeader := range signedHeaders {
		header := strings.ToLower(signedHeader)

		if header == "host" {
			canonicalRequest.WriteString("host:")
			canonicalRequest.WriteString(strings.TrimSpace(r.Host))
			canonicalRequest.WriteString("\n")
			continue
		}

		canonicalRequest.WriteString(header)
		canonicalRequest.WriteString(":")
		canonicalRequest.WriteString(strings.TrimSpace(r.Header.Get(signedHeader)))
		canonicalRequest.WriteString("\n")
	}

	canonicalRequest.WriteString("\n")

	canonicalRequest.WriteString(strings.Join(signedHeaders, ";"))
	canonicalRequest.WriteString("\n")

	var recomputedPayloadHash string

	if payloadHash == UnsignedPayload || payloadHash == StreamingPayload {
		recomputedPayloadHash = payloadHash
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return "", errors.New("could not read body")
		}

		slog.Debug("Read body", "length", len(body))

		// Restore body
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		hash := sha256.Sum256(body)
		recomputedPayloadHash = hex.EncodeToString(hash[:])

		slog.Debug("Recomputed payload hash", "sha256", security.Trunc(recomputedPayloadHash))
	}

	canonicalRequest.WriteString(recomputedPayloadHash)

	return canonicalRequest.String(), nil
}

func buildCanonicalURI(path string) string {
	parts := strings.Split(path, "/")

	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}

	canonicalURI := strings.Join(parts, "/")

	return canonicalURI
}

func buildQueryString(rawQuery string) string {
	m, _ := url.ParseQuery(rawQuery)

	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var parts []string

	for _, key := range keys {
		values := m[key]
		sort.Strings(values)

		for _, value := range values {
			parts = append(parts, queryEncode(key)+"="+queryEncode(value))
		}
	}

	return strings.Join(parts, "&")
}

func queryEncode(kv string) string {
	esc := url.QueryEscape(kv)
	esc = strings.ReplaceAll(esc, "+", "%20")
	esc = strings.ReplaceAll(esc, "%7E", "~")
	return esc
}

func buildStringToSign(
	timestamp string,
	scope string,
	canonicalRequest string,
) string {
	var stringToSign strings.Builder

	stringToSign.WriteString("AWS4-HMAC-SHA256")
	stringToSign.WriteString("\n")

	stringToSign.WriteString(timestamp)
	stringToSign.WriteString("\n")

	stringToSign.WriteString(scope)
	stringToSign.WriteString("\n")

	hash := sha256.Sum256([]byte(canonicalRequest))
	stringToSign.WriteString(hex.EncodeToString(hash[:]))

	return stringToSign.String()
}

func computeSignature(
	secretKey string,
	scope string,
	stringToSign string,
) (string, error) {
	scopeParts := strings.Split(scope, "/")

	if len(scopeParts) != 4 {
		return "", errors.New("scope must contain 4 parts")
	}

	date := scopeParts[0]
	region := scopeParts[1]
	service := scopeParts[2]

	dateKey := hmacSHA256([]byte("AWS4"+secretKey), []byte(date))
	dateRegionKey := hmacSHA256(dateKey, []byte(region))
	dateRegionServiceKey := hmacSHA256(dateRegionKey, []byte(service))
	signingKey := hmacSHA256(dateRegionServiceKey, []byte("aws4_request"))
	signature := hmacSHA256(signingKey, []byte(stringToSign))
	signatureString := hex.EncodeToString(signature)

	return signatureString, nil
}

func hmacSHA256(key, value []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(value)
	return mac.Sum(nil)
}
