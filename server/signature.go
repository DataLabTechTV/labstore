package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func buildCanonicalRequest(r *http.Request, signedHeaders []string) (string, bool) {
	var canonicalRequest strings.Builder

	canonicalRequest.WriteString(r.Method)
	canonicalRequest.WriteString("\n")

	canonicalRequest.WriteString(r.URL.Path)
	canonicalRequest.WriteString("\n")

	canonicalRequest.WriteString(r.URL.RawQuery)
	canonicalRequest.WriteString("\n")

	for _, signedHeader := range signedHeaders {
		header := strings.ToLower(signedHeader)

		if header == "host" {
			canonicalRequest.WriteString("host:")
			canonicalRequest.WriteString(r.Host)
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
		return "", false
	}

	canonicalRequest.WriteString(fmt.Sprintf("%x", sha256.Sum256(body)))

	return canonicalRequest.String(), true
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
) (string, bool) {
	scopeParts := strings.Split(scope, "/")
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

	return string(signature), true
}

func verifyAWSSigV4(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")

	log.Debug("Authorization: " + auth)

	// Remove prefix

	if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256") {
		return "", false
	}

	auth, ok := strings.CutPrefix(auth, "AWS4-HMAC-SHA256 ")
	if !ok {
		return "", false
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
		return "", false
	}

	if len(signedHeaders) == 0 {
		return "", false
	}

	if signature == "" {
		return "", false
	}

	log.Debug("Credentials: " + credentials)
	log.Debug("SignedHeaders: " + strings.Join(signedHeaders, ";"))
	log.Debug("Signature: " + signature)

	// Extract access key and scope

	credentialParts := strings.Split(credentials, "/")

	accessKey := credentialParts[0]
	secretKey, ok := users[accessKey]

	if !ok {
		return "", false
	}

	scope := strings.Join(credentialParts[1:], "/")

	log.Debug("Access key: " + accessKey)
	log.Debug("Scope: " + scope)

	// Compute signature

	canonicalRequest, ok := buildCanonicalRequest(r, signedHeaders)

	if !ok {
		return "", false
	}

	log.Debug("Canonical request: " + canonicalRequest)

	timestamp := r.Header.Get("X-Amz-Date")
	stringToSign := buildStringToSign(timestamp, scope, canonicalRequest)

	log.Debug("String to sign: " + stringToSign)

	recomputedSignature, ok := computeSignature(secretKey, scope, stringToSign)

	if !ok {
		return "", false
	}

	log.Debug("Signature (recomputed): " + fmt.Sprintf("%x", recomputedSignature))

	if signature == recomputedSignature {
		return accessKey, true
	}

	log.Error("Original and recomputed signatures differ")

	return "", false
}
