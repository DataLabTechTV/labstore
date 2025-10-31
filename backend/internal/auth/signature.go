package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/DataLabTechTV/labstore/backend/pkg/iam"
	"github.com/DataLabTechTV/labstore/backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

const unsignedPayload = "UNSIGNED-PAYLOAD"

func VerifyAWSSigV4(c *fiber.Ctx) (string, error) {
	auth := c.Get("Authorization")
	logger.Log.Debug("Authorization: " + auth)

	payloadHash := c.Get("X-Amz-Content-SHA256")
	logger.Log.Debug("X-Amz-Content-SHA256: " + payloadHash)

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

	logger.Log.Debug("Credentials: " + credentials)
	logger.Log.Debug("SignedHeaders: " + strings.Join(signedHeaders, ";"))
	logger.Log.Debug("Signature: " + signature)

	// Extract access key and scope

	credentialParts := strings.Split(credentials, "/")

	accessKey := credentialParts[0]
	logger.Log.Debug("Access key: " + accessKey)

	secretKey, ok := iam.Users[accessKey]
	if !ok {
		return "", fmt.Errorf("no secret key found for access key %s", accessKey)
	}

	scope := strings.Join(credentialParts[1:], "/")
	logger.Log.Debug("Scope: " + scope)

	// Compute signature

	canonicalRequest, err := buildCanonicalRequest(c, signedHeaders, payloadHash)
	if err != nil {
		return "", errors.New("could not build canonical request")
	}
	logger.Log.Debug("Canonical request: " + canonicalRequest)

	timestamp := c.Get("X-Amz-Date")
	logger.Log.Debug("Timestamp: " + timestamp)

	stringToSign := buildStringToSign(timestamp, scope, canonicalRequest)
	logger.Log.Debug("String to sign: " + stringToSign)

	recomputedSignature, err := computeSignature(secretKey, scope, stringToSign)
	if err != nil {
		return "", errors.New("could not compute signature")
	}

	logger.Log.Debug("Signature (recomputed): " + recomputedSignature)

	byteSignature, err := hex.DecodeString(signature)
	if err != nil {
		return "", errors.New("could not decode original signature")
	}

	byteRecomputedSignature, err := hex.DecodeString((recomputedSignature))
	if err != nil {
		return "", errors.New("could not decode recomputed signature")
	}

	if hmac.Equal(byteSignature, byteRecomputedSignature) {
		return accessKey, nil
	}

	logger.Log.Error("Original and recomputed signatures differ")
	return "", errors.New("signatures do not match")
}

func buildCanonicalRequest(
	c *fiber.Ctx,
	signedHeaders []string,
	payloadHash string,
) (string, error) {
	var canonicalRequest strings.Builder

	canonicalRequest.WriteString(c.Method())
	canonicalRequest.WriteString("\n")

	canonicalRequest.WriteString(c.Path())
	canonicalRequest.WriteString("\n")

	queryString := buildQueryString(string(c.Request().URI().QueryString()))
	logger.Log.Debug("Canonical query string: " + queryString)
	canonicalRequest.WriteString(queryString)
	canonicalRequest.WriteString("\n")

	for _, signedHeader := range signedHeaders {
		header := strings.ToLower(signedHeader)

		if header == "host" {
			canonicalRequest.WriteString("host:")
			canonicalRequest.WriteString(strings.TrimSpace(c.Hostname()))
			canonicalRequest.WriteString("\n")
			continue
		}

		canonicalRequest.WriteString(header)
		canonicalRequest.WriteString(":")
		canonicalRequest.WriteString(strings.TrimSpace(c.Get(signedHeader)))
		canonicalRequest.WriteString("\n")
	}

	canonicalRequest.WriteString("\n")

	canonicalRequest.WriteString(strings.Join(signedHeaders, ";"))
	canonicalRequest.WriteString("\n")

	var recomputedPayloadHash string

	if payloadHash == unsignedPayload {
		recomputedPayloadHash = payloadHash
	} else {
		body := c.Body()
		logger.Log.Debugf("Body length: %d", len(body))

		// Restore body
		c.Request().SetBody(body)

		hash := sha256.Sum256(body)
		recomputedPayloadHash = hex.EncodeToString(hash[:])
	}

	canonicalRequest.WriteString(recomputedPayloadHash)

	return canonicalRequest.String(), nil
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

	return hex.EncodeToString(signature), nil
}

func hmacSHA256(key, value []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(value)
	return mac.Sum(nil)
}
