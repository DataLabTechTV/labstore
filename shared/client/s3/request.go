package s3

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/auth"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/security"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

const DefaultRequestTimeout = 1 * time.Minute
const DefaultRegion = "eu-west-1"

func (client *S3Client) DoSigV4Request(method, rawURL string, body io.ReadCloser) (*http.Response, error) {
	r, err := http.NewRequestWithContext(client.Ctx, method, rawURL, body)
	if err != nil {
		return nil, err
	}

	credential := &auth.SigV4Credential{
		AccessKey: client.AccessKey,
		SecretKey: client.SecretKey,
		Scope:     buildScope(),
	}

	authorization := &auth.SigV4Authorization{
		Credential: credential,
		SignedHeaders: []string{
			"host",
			"x-amz-content-sha256",
			"x-amz-date",
		},
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	var payloadHash string

	switch r.Body.(type) {
	case *SigV4ChunkEncoder:
		payloadHash = auth.StreamingPayload
	default:
		payloadHash, r.Body, err = computePayloadHash(r.Body)
		if err != nil {
			return nil, err
		}
	}

	timestamp := time.Now().Format(types.CompactISO8601)

	r.Header.Set("X-Amz-Content-SHA256", payloadHash)
	r.Header.Set("X-Amz-Date", timestamp)

	sigV4Req := &auth.SigV4Request{
		Method:               method,
		CanonicalURI:         auth.BuildCanonicalURI(parsedURL.Path),
		CanonicalQueryString: auth.BuildCanonicalQueryString(parsedURL.RawQuery),
		CanonicalHeaders:     auth.BuildCanonicalHeaders(r, authorization),
		Authorization:        authorization,
		Timestamp:            timestamp,
		PayloadHash:          payloadHash,
	}

	stringToSign := sigV4Req.BuildStringToSign()
	slog.Debug("string to sign", "string_to_sign", security.TruncLastLine(stringToSign))

	authorization.Signature, err = auth.ComputeSignature(sigV4Req.Authorization.Credential, stringToSign)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Authorization", authorization.String())

	switch enc := r.Body.(type) {
	case *SigV4ChunkEncoder:
		enc.chunk.Ctx.Signature = sigV4Req.Authorization.Signature
		enc.chunk.Ctx.Credential = sigV4Req.Authorization.Credential
		enc.chunk.Ctx.Timestamp = sigV4Req.Timestamp
		enc.chunk.Ctx.IsStreaming = true
		r.Header.Set("Content-Type", "application/octet")
		r.Header.Set("Content-Length", fmt.Sprint(enc.TotalSize))
		r.Header.Set("X-Amz-Decoded-Content-Length", fmt.Sprint(enc.DataSize))
		r.Body = io.NopCloser(bufio.NewReaderSize(r.Body, config.S3.IO.BufferSize))
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func buildScope() string {
	var scope strings.Builder

	scope.WriteString(time.Now().Format("20060102"))
	scope.WriteRune('/')

	scope.WriteString(DefaultRegion)
	scope.WriteRune('/')

	scope.WriteString("s3/aws4_request")

	return scope.String()
}

func computePayloadHash(body io.ReadCloser) (string, io.ReadCloser, error) {
	if body == nil {
		bytePayloadHash := sha256.Sum256([]byte{})
		payloadHash := hex.EncodeToString(bytePayloadHash[:])
		return payloadHash, nil, nil
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return "", nil, err
	}

	newBody := io.NopCloser(bytes.NewBuffer(data))

	bytePayloadHash := sha256.Sum256(data)
	payloadHash := hex.EncodeToString(bytePayloadHash[:])

	return payloadHash, newBody, nil
}
