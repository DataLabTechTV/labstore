package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

const (
	UnsignedPayload  = "UNSIGNED-PAYLOAD"
	StreamingPayload = "STREAMING-AWS4-HMAC-SHA256-PAYLOAD"
)

func ComputeSignature(cred *SigV4Credential, stringToSign string) (string, error) {
	scopeParts := strings.Split(cred.Scope, "/")

	if len(scopeParts) != 4 {
		return "", errors.New("scope must contain 4 parts")
	}

	date := scopeParts[0]
	region := scopeParts[1]
	service := scopeParts[2]

	dateKey := hmacSHA256([]byte("AWS4"+cred.SecretKey), []byte(date))
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
