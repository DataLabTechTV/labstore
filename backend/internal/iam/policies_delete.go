package iam

import (
	"context"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

type DeletePolicyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DeletePolicyResponse"`
	ResponseMetadata *ResponseMetadata
}

func (store *Store) DeletePolicy(ctx context.Context, arn string) error {
	policy, err := store.GetPolicyByArn(ctx, arn)
	if err != nil {
		slog.Error("delete policy", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypePolicy, Resource: arn}
	}

	query := `
	DELETE FROM policies
	WHERE arn = $1
	`

	_, err = store.sqlExecContext(ctx, query, policy.Arn)
	if err != nil {
		slog.Error("delete policy", "err", err)
		return err
	}

	delete(store.CachedPolicies, policy.PolicyID)

	return nil
}

func DeletePolicyHandler(w http.ResponseWriter, r *http.Request) {
	policyArn := r.Form.Get("PolicyArn")
	if policyArn == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("PolicyArn"))
		return
	}

	ctx := r.Context()

	if err := store.DeletePolicy(ctx, policyArn); err != nil {
		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMNoSuchEntity(string(errNotFound.Type), errNotFound.Resource))
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	response := &DeletePolicyResponse{
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
