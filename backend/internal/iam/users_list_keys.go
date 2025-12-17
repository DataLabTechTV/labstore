package iam

import (
	"encoding/xml"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

type ListAccessKeysResponse struct {
	XMLName              xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ ListAccessKeysResponse"`
	ListAccessKeysResult *ListAccessKeysResult
	ResponseMetadata     *ResponseMetadata
}

type ListAccessKeysResult struct {
	UserName          string
	AccessKeyMetadata *AccessKeyMetadata
	IsTruncated       bool
}

type AccessKeyMetadata struct {
	Member []*AccessKeyMetadataMember `xml:"member"`
}

type AccessKeyMetadataMember struct {
	UserName    string
	AccessKeyId string
	Status      AccessKeyStatus
}

func ListAccessKeysHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	user, err := store.GetUserByName(userName)
	if err != nil {
		errs.Handle(w, errs.IAMNoSuchEntity(string(errs.ErrEntityTypeUser), userName))
		return
	}

	members := []*AccessKeyMetadataMember{}
	if user.AccessKeyID.Valid {
		members = append(members, &AccessKeyMetadataMember{
			UserName:    user.Name,
			AccessKeyId: user.AccessKeyID.String,
			Status:      AccessKeyActive,
		})
	}

	response := &ListAccessKeysResponse{
		ListAccessKeysResult: &ListAccessKeysResult{
			UserName: user.Name,
			AccessKeyMetadata: &AccessKeyMetadata{
				Member: members,
			},
		},
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
