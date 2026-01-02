package iam

import (
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

func ListAccessKeysHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.Form.Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	ctx := r.Context()

	user, err := store.GetUserByName(ctx, userName)
	if err != nil {
		errs.Handle(w, errs.IAMNoSuchEntity(string(errs.ErrEntityTypeUser), userName))
		return
	}

	members := []*t.AccessKeyMetadataMember{}
	if user.AccessKeyID.Valid {
		members = append(members, &t.AccessKeyMetadataMember{
			UserName:    user.Name,
			AccessKeyId: user.AccessKeyID.String,
			Status:      t.AccessKeyActive,
		})
	}

	response := &t.ListAccessKeysResponse{
		ListAccessKeysResult: &t.ListAccessKeysResult{
			UserName: user.Name,
			AccessKeyMetadata: &t.AccessKeyMetadata{
				Member: members,
			},
			IsTruncated: false,
		},
		ResponseMetadata: &t.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	helper.WriteXMLResponse(w, http.StatusOK, response)
}
