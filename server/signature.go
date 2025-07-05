package server

import (
	"net/http"
	"strings"
)

// --- AWS Signature V4 Verification ---

// For simplicity, we only check Authorization header with access key and sign
// string presence. Full signature verification is complex, so this is a stub.

func verifyAWSSigV4(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	// Example: Authorization: AWS4-HMAC-SHA256 Credential=myaccesskey/20230629/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-date, Signature=...
	if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256") {
		return "", false
	}
	// Parse Credential= part
	parts := strings.Split(auth, ",")
	var credential string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if after, ok := strings.CutPrefix(p, "Credential="); ok {
			credential = after
			break
		}
	}
	if credential == "" {
		return "", false
	}
	accessKey := strings.Split(credential, "/")[0]
	_, ok := users[accessKey]
	if !ok {
		return "", false
	}
	// Here you'd do full signature verification with secretKey, canonical request, etc.
	// For MVP, we just check presence of accessKey in users.
	return accessKey, true
}
