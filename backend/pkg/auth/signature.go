package auth

import (
	"bytes"
	"context"
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

	"github.com/IllumiKnowLabs/labstore/backend/internal/iam"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/security"
)

type SigV4Request struct {
	Method               string
	CanonicalURI         string
	CanonicalQueryString string
	CanonicalHeaders     map[string]string
	Authorization        *SigV4Authorization
	Timestamp            string
	PayloadHash          string
}

type SigV4Authorization struct {
	Credential    *SigV4Credential
	SignedHeaders []string
	Signature     string
}

type SigV4Credential struct {
	AccessKey string
	SecretKey string
	Scope     string
}

type SigV4Context struct {
	Signature   string
	Credential  *SigV4Credential
	Timestamp   string
	IsStreaming bool
}

func VerifySigV4(r *http.Request) (*SigV4Context, error) {
	req, err := newSigV4Request(r)
	if err != nil {
		return nil, fmt.Errorf("sigv4: %w", err)
	}

	if err := req.validatePayloadHash(r); err != nil {
		return nil, fmt.Errorf("sigv4: %w", err)
	}

	res, err := req.validateSignature()
	if err != nil {
		return nil, fmt.Errorf("sigv4: %w", err)
	}

	return res, nil
}

func newSigV4Request(r *http.Request) (*SigV4Request, error) {
	authorization := r.Header.Get("Authorization")
	slog.Debug("parsing sigv4 request", "authorization", security.TruncParamHeader(authorization, "Signature"))

	auth, err := newSigV4Authorization(authorization)
	if err != nil {
		return nil, err
	}

	payloadHash := r.Header.Get("X-Amz-Content-SHA256")
	slog.Debug("payload hash", "x-amz-content-sha256", security.Trunc(payloadHash))

	timestamp := r.Header.Get("X-Amz-Date")
	slog.Debug("timestamp", "x-amz-date", timestamp)

	canonicalURI := BuildCanonicalURI(r.URL.Path)
	slog.Debug("canonical uri", "uri", canonicalURI)

	canonicalQueryString := BuildCanonicalQueryString(r.URL.RawQuery)
	slog.Debug("canonical query string", "query_string", canonicalQueryString)

	canonicalHeaders := BuildCanonicalHeaders(r, auth)
	slog.Debug("canonical headers", "headers", canonicalHeaders)

	res := &SigV4Request{
		Method:               r.Method,
		CanonicalURI:         canonicalURI,
		CanonicalQueryString: canonicalQueryString,
		CanonicalHeaders:     canonicalHeaders,
		Authorization:        auth,
		Timestamp:            timestamp,
		PayloadHash:          payloadHash,
	}

	return res, nil
}

// Check for SigV4 prefix, and extract credential, signed headers and signature
func newSigV4Authorization(authorization string) (*SigV4Authorization, error) {
	auth, ok := strings.CutPrefix(authorization, "AWS4-HMAC-SHA256 ")
	if !ok {
		return nil, errors.New("header Authorization must start with AWS4-HMAC-SHA256")
	}

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
		"authorization",
		"credential", credential,
		"signed_headers", strings.Join(signedHeaders, ";"),
		"signature", security.Trunc(signature),
	)

	cred, err := newSigV4Credential(credential)
	if err != nil {
		return nil, err
	}

	res := &SigV4Authorization{
		Credential:    cred,
		SignedHeaders: signedHeaders,
		Signature:     signature,
	}

	return res, nil
}

// Extract access key and scope, and retrieve secret key from IAM
func newSigV4Credential(credential string) (*SigV4Credential, error) {
	credentialParts := strings.Split(credential, "/")

	accessKey := credentialParts[0]

	ctx := context.Background()

	user, err := iam.GetStore().GetUserByAccessKey(ctx, accessKey)
	if err != nil || !user.AccessKeyID.Valid {
		slog.Error("get user by access key", "err", err)
		return nil, errors.New("invalid access key")
	}

	scope := strings.Join(credentialParts[1:], "/")

	slog.Debug("credential", "access_key", accessKey, "scope", scope)

	plainSecretKey, err := security.DecryptAESGCM(user.SecretKey, config.Storage.MasterKeyPath)
	if err != nil {
		return nil, err
	}

	res := &SigV4Credential{
		AccessKey: user.AccessKeyID.String,
		SecretKey: plainSecretKey,
		Scope:     scope,
	}

	return res, nil
}

func BuildCanonicalURI(path string) string {
	parts := strings.Split(path, "/")

	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}

	canonicalURI := strings.Join(parts, "/")

	return canonicalURI
}

func BuildCanonicalQueryString(rawQuery string) string {
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

func BuildCanonicalHeaders(r *http.Request, auth *SigV4Authorization) map[string]string {
	headers := make(map[string]string)

	for _, signedHeader := range auth.SignedHeaders {
		header := strings.ToLower(signedHeader)

		var value string

		if header == "host" {
			value = r.Host
		} else {
			value = r.Header.Get(signedHeader)
		}

		headers[header] = strings.TrimSpace(value)
	}

	return headers
}

func (req *SigV4Request) validatePayloadHash(r *http.Request) error {
	if req.PayloadHash == UnsignedPayload || req.PayloadHash == StreamingPayload {
		return nil
	}

	// Recompute body hash
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("could not read body")
	}

	slog.Debug("body", "length", len(body))

	// Restore body
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	bytePayloadHash, err := hex.DecodeString(req.PayloadHash)
	if err != nil {
		return errors.New("could not decode payload hash")
	}

	byteRecomputedPayloadHash := sha256.Sum256(body)
	recomputedPayloadHash := hex.EncodeToString(byteRecomputedPayloadHash[:])

	slog.Debug(
		"comparing payload hashes",
		"received", security.Trunc(req.PayloadHash),
		"recomputed", security.Trunc(recomputedPayloadHash),
	)

	if hmac.Equal(bytePayloadHash, byteRecomputedPayloadHash[:]) {
		return nil
	}

	slog.Error("payload hashes differ")
	return errors.New("payload hashes do not match")
}

// Recompute and validate SigV4 signature
func (req *SigV4Request) validateSignature() (*SigV4Context, error) {
	stringToSign := req.BuildStringToSign()
	slog.Debug("string to sign", "string_to_sign", security.TruncLastLine(stringToSign))

	signature, err := ComputeSignature(req.Authorization.Credential, stringToSign)
	if err != nil {
		return nil, errors.New("could not compute signature")
	}

	byteSignature, err := hex.DecodeString(req.Authorization.Signature)
	if err != nil {
		return nil, errors.New("could not decode original signature")
	}

	byteRecomputedSignature, err := hex.DecodeString((signature))
	if err != nil {
		return nil, errors.New("could not decode recomputed signature")
	}

	slog.Debug(
		"comparing signatures",
		"received", security.Trunc(req.Authorization.Signature),
		"recomputed", security.Trunc(signature),
	)

	if hmac.Equal(byteSignature, byteRecomputedSignature) {
		isStreaming := req.PayloadHash == StreamingPayload

		res := &SigV4Context{
			Credential:  req.Authorization.Credential,
			Signature:   req.Authorization.Signature,
			Timestamp:   req.Timestamp,
			IsStreaming: isStreaming,
		}

		return res, nil
	}

	slog.Error("signatures differ")
	return nil, errors.New("signatures do not match")
}

func (req *SigV4Request) buildCanonicalRequest() string {
	var canonicalRequest strings.Builder

	canonicalRequest.WriteString(req.Method)
	canonicalRequest.WriteString("\n")

	canonicalRequest.WriteString(req.CanonicalURI)
	canonicalRequest.WriteString("\n")

	canonicalRequest.WriteString(req.CanonicalQueryString)
	canonicalRequest.WriteString("\n")

	for _, header := range req.Authorization.SignedHeaders {
		canonicalRequest.WriteString(header)
		canonicalRequest.WriteString(":")
		canonicalRequest.WriteString(req.CanonicalHeaders[header])
		canonicalRequest.WriteString("\n")
	}

	canonicalRequest.WriteString("\n")

	canonicalRequest.WriteString(strings.Join(req.Authorization.SignedHeaders, ";"))
	canonicalRequest.WriteString("\n")

	canonicalRequest.WriteString(req.PayloadHash)

	return canonicalRequest.String()
}

func (req *SigV4Request) BuildStringToSign() string {
	canonicalRequest := req.buildCanonicalRequest()
	slog.Debug("canonical request", "canonical_request", security.TruncLastLine(canonicalRequest))

	var stringToSign strings.Builder

	stringToSign.WriteString("AWS4-HMAC-SHA256")
	stringToSign.WriteString("\n")

	stringToSign.WriteString(req.Timestamp)
	stringToSign.WriteString("\n")

	stringToSign.WriteString(req.Authorization.Credential.Scope)
	stringToSign.WriteString("\n")

	hash := sha256.Sum256([]byte(canonicalRequest))
	stringToSign.WriteString(hex.EncodeToString(hash[:]))

	return stringToSign.String()
}

func (auth *SigV4Authorization) String() string {
	var b strings.Builder

	b.WriteString("AWS4-HMAC-SHA256 Credential=")
	b.WriteString(auth.Credential.AccessKey)
	b.WriteRune('/')
	b.WriteString(auth.Credential.Scope)
	b.WriteString(", SignedHeaders=")
	b.WriteString(strings.Join(auth.SignedHeaders, ";"))
	b.WriteString(", Signature=")
	b.WriteString(auth.Signature)

	return b.String()
}
