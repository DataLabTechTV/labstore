package handlers

import (
	"github.com/IllumiKnowLabs/labstore/client/pkg/iam"
)

type IAMHandler struct {
	Client *iam.Client
}

func NewIAMHandler(client *iam.Client) *IAMHandler {
	return &IAMHandler{
		Client: client,
	}
}
