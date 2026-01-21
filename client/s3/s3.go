package s3

import (
	"context"
	"fmt"
	"net/url"
)

type Client struct {
	Ctx       context.Context
	Host      string
	Port      uint16
	AccessKey string
	SecretKey string
	TLS       bool
	ChunkSize int

	baseURL *url.URL
}

func NewS3Client(
	ctx context.Context,
	host string,
	port uint16,
	accessKey string,
	secretKey string,
	tls bool,
	chunkSize int,
) *Client {
	client := &Client{
		Ctx:       ctx,
		Host:      host,
		Port:      port,
		AccessKey: accessKey,
		SecretKey: secretKey,
		TLS:       tls,
		ChunkSize: chunkSize,
	}

	var scheme string
	if tls {
		scheme = "http"
	} else {
		scheme = "http"
	}

	client.baseURL = &url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%d", client.Host, client.Port),
	}

	return client
}

func (c *Client) IsDone() bool {
	select {
	case <-c.Ctx.Done():
		return true
	default:
		return false
	}
}
