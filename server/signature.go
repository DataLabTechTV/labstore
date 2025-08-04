package server

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func toCanonicalRequest(r *http.Request, timestamp string) (string, bool) {
	var canReq strings.Builder

	canReq.WriteString(r.Method)
	canReq.WriteString("\n")

	canReq.WriteString(r.URL.Path)
	canReq.WriteString("\n")

	canReq.WriteString(r.URL.RawQuery)
	canReq.WriteString("\n")

	canReq.WriteString("host:")
	canReq.WriteString(r.Host)
	canReq.WriteString("\n")

	canReq.WriteString("x-amz-date:")
	canReq.WriteString(timestamp)
	canReq.WriteString("\n\n")

	canReq.WriteString("host;x-amz-date")
	canReq.WriteString("\n")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		return "", false
	}

	canReq.WriteString(fmt.Sprintf("%x", sha256.Sum256(body)))

	return canReq.String(), true
}

func verifyAWSSigV4(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")

	if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256") {
		return "", false
	}

	var ok bool

	auth, ok = strings.CutPrefix(auth, "AWS4-HMAC-SHA256 ")
	if !ok {
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
	_, ok = users[accessKey]
	if !ok {
		return "", false
	}

	// Here you'd do full signature verification with secretKey, canonical request, etc.
	// For MVP, we just check presence of accessKey in users.

	timestamp := r.Header.Get("X-Amz-Date")
	// canonicalRequest, ok := toCanonicalRequest(r, timestamp)
	_, ok = toCanonicalRequest(r, timestamp)

	if !ok {
		return "", false
	}

	return accessKey, true
}
