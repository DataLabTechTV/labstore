package iam

import (
	"net/http"

	"github.com/IllumiKnowLabs/labstore/server/internal/core"
	"github.com/IllumiKnowLabs/labstore/server/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/types"
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

	members := []*types.AccessKeyMetadataMember{}
	if user.AccessKeyID.Valid {
		members = append(members, &types.AccessKeyMetadataMember{
			UserName:    user.Name,
			AccessKeyId: user.AccessKeyID.String,
			Status:      types.AccessKeyActive,
		})
	}

	response := &types.ListAccessKeysResponse{
		ListAccessKeysResult: &types.ListAccessKeysResult{
			UserName: user.Name,
			AccessKeyMetadata: &types.AccessKeyMetadata{
				Member: members,
			},
			IsTruncated: false,
		},
		ResponseMetadata: &types.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	helper.WriteXMLResponse(w, http.StatusOK, response)
}
