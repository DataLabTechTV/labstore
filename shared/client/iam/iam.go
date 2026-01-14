package iam

import (
	"context"
	"fmt"
	"net/url"
)

type Client struct {
	Ctx  context.Context
	Host string
	Port uint16

	baseURL *url.URL
}

func NewIAMClient(ctx context.Context, host string, port uint16) *Client {
	client := &Client{
		Ctx:  ctx,
		Host: host,
		Port: port,
	}

	client.baseURL = &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", client.Host, client.Port),
	}

	return client
}
