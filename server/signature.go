package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

func queryEncode(kv string) string {
	esc := url.QueryEscape(kv)
	esc = strings.ReplaceAll(esc, "+", "%20")
	esc = strings.ReplaceAll(esc, "%7E", "~")
	return esc
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

func buildCanonicalRequest(r *http.Request, signedHeaders []string) (string, error) {
	var canonicalRequest strings.Builder

	canonicalRequest.WriteString(r.Method)
	canonicalRequest.WriteString("\n")

	canonicalRequest.WriteString(r.URL.Path)
	canonicalRequest.WriteString("\n")

	queryString := buildQueryString(r.URL.RawQuery)
	log.Debug("Canonical query string: " + queryString)
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

	body, err := io.ReadAll(r.Body)

	if err != nil {
		return "", errors.New("could not read body")
	}

	if len(body) == 0 {
		log.Debug(fmt.Sprintf("Body (%d bytes): EMPTY", len(body)))
	} else {
		log.Debug(fmt.Sprintf("Body (%d bytes): %s", len(body), string(body)))
	}

	// Restore body
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	hash := sha256.Sum256(body)
	canonicalRequest.WriteString(fmt.Sprintf("%x", hash))

	return canonicalRequest.String(), nil
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

	stringToSign.WriteString(fmt.Sprintf("%x", sha256.Sum256([]byte(canonicalRequest))))
	stringToSign.WriteString("\n")

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
	awsRegion := scopeParts[1]
	awsService := scopeParts[2]

	hashFunc := sha256.New

	mac := hmac.New(hashFunc, []byte("AWS4"+secretKey))
	mac.Write([]byte(date))
	dateKey := mac.Sum(nil)

	mac = hmac.New(hashFunc, dateKey)
	mac.Write([]byte(awsRegion))
	dateRegionKey := mac.Sum(nil)

	mac = hmac.New(hashFunc, dateRegionKey)
	mac.Write([]byte(awsService))
	dateRegionServiceKey := mac.Sum(nil)

	mac = hmac.New(hashFunc, dateRegionServiceKey)
	mac.Write([]byte("aws4_request"))
	signingKey := mac.Sum(nil)

	mac = hmac.New(hashFunc, signingKey)
	mac.Write([]byte(stringToSign))
	signature := mac.Sum(nil)

	return hex.EncodeToString(signature), nil
}

func verifyAWSSigV4(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")

	log.Debug("Authorization: " + auth)

	// Remove prefix

	if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256") {
		return "", errors.New("header Authorization must start with AWS4-HMAC-SHA256")
	}

	auth, ok := strings.CutPrefix(auth, "AWS4-HMAC-SHA256 ")
	if !ok {
		return "", errors.New("could not remove prefix AWS4-HMAC-SHA256")
	}

	// Parse credentials, signed headers, and signature

	parts := strings.Split(auth, ",")

	var credentials string
	var signedHeaders []string
	var signature string

	for _, p := range parts {
		p = strings.TrimSpace(p)

		if after, ok := strings.CutPrefix(p, "Credential="); ok {
			credentials = after
		}

		if after, ok := strings.CutPrefix(p, "SignedHeaders="); ok {
			signedHeaders = strings.Split(after, ";")
		}

		if after, ok := strings.CutPrefix(p, "Signature="); ok {
			signature = after
		}
	}

	if credentials == "" {
		return "", errors.New("header Credentials is empty")
	}

	if len(signedHeaders) == 0 {
		return "", errors.New("header SignedHeaders is empty")
	}

	if signature == "" {
		return "", errors.New("header Signature is empty")
	}

	log.Debug("Credentials: " + credentials)
	log.Debug("SignedHeaders: " + strings.Join(signedHeaders, ";"))
	log.Debug("Signature: " + signature)

	// Extract access key and scope

	credentialParts := strings.Split(credentials, "/")

	accessKey := credentialParts[0]
	secretKey, ok := users[accessKey]

	if !ok {
		return "", fmt.Errorf("No secret key found for access key %s", accessKey)
	}

	scope := strings.Join(credentialParts[1:], "/")

	log.Debug("Access key: " + accessKey)
	log.Debug("Scope: " + scope)

	// Compute signature

	canonicalRequest, err := buildCanonicalRequest(r, signedHeaders)

	if err != nil {
		return "", errors.New("could not build canonical request")
	}

	log.Debug("Canonical request: " + canonicalRequest)

	timestamp := r.Header.Get("X-Amz-Date")
	stringToSign := buildStringToSign(timestamp, scope, canonicalRequest)

	log.Debug("String to sign: " + stringToSign)

	recomputedSignature, err := computeSignature(secretKey, scope, stringToSign)

	if err != nil {
		return "", errors.New("could not compute signature")
	}

	log.Debug("Signature (recomputed): " + recomputedSignature)

	if signature == recomputedSignature {
		return accessKey, nil
	}

	log.Error("Original and recomputed signatures differ")

	return "", errors.New("signatures do not match")
}
