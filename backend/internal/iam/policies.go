package iam

import (
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

const (
	defaultPolicyPath    = "/"
	defaultPolicyVersion = "v1"

	latestPolicyDocumentVersion = "2012-10-17"

	adminPolicy = "admin-policy"
)

type Policy struct {
	PolicyID string `db:"policy_id"`
	Name     string `db:"name"`
	Arn      string `db:"arn"`

	Document *PolicyDocument `db:"document"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	AttachmentCount uint
}

type PolicyDocument struct {
	Version   string
	Statement []Statement
}

type Statement struct {
	Effect   Effect
	Action   Actions
	Resource Resources
}

type CreatePolicyResponse struct {
	XMLName            xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ CreatePolicyResponse"`
	CreatePolicyResult *CreatePolicyResult
	ResponseMetadata   *ResponseMetadata
}

type CreatePolicyResult struct {
	Policy *PolicyResult
}

type PolicyResult struct {
	XMLName          xml.Name `xml:"Policy"`
	PolicyName       string
	DefaultVersionId string
	PolicyId         string
	Path             string
	Arn              string
	AttachmentCount  uint
	CreateDate       time.Time
	UpdateDate       time.Time
}

func (pd *PolicyDocument) Value() (driver.Value, error) {
	return json.Marshal(pd)
}

func (pd *PolicyDocument) Scan(src any) error {
	if src == nil {
		*pd = PolicyDocument{}
		return nil
	}

	switch s := src.(type) {
	case []byte:
		return json.Unmarshal(s, pd)
	case string:
		return json.Unmarshal([]byte(s), pd)
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
}

func GetPolicyByID(policyID string) (*Policy, error) {
	var policy Policy

	query := `SELECT * FROM policies WHERE policy_id = $1`
	if err := store.readDB.Get(&policy, query, policyID); err != nil {
		return nil, err
	}

	var attachments int
	query = `
	SELECT count(*)
	FROM (
		SELECT 1 FROM user_policies WHERE policy_id = $1
		UNION ALL
		SELECT 1 FROM group_policies WHERE policy_id = $1
	)
	`

	if err := store.readDB.Get(&attachments, query, policyID); err != nil {
		return nil, err
	}

	return &policy, nil
}

func GetPolicyByArn(arn string) (*Policy, error) {
	var policy Policy

	query := `SELECT * FROM policies WHERE arn = $1`
	if err := store.readDB.Get(&policy, query, arn); err != nil {
		return nil, err
	}

	var attachments int
	query = `
	SELECT count(*)
	FROM (
		SELECT 1 FROM user_policies WHERE policy_id = $1
		UNION ALL
		SELECT 1 FROM group_policies WHERE policy_id = $1
	)
	`

	if err := store.readDB.Get(&attachments, query, policy.PolicyID); err != nil {
		return nil, err
	}

	return &policy, nil
}

func CreatePolicy(name string, doc *PolicyDocument) (*Policy, error) {
	policyID := GenerateUniqueID(ManagedPolicyUniqueID)

	policy := &Policy{
		PolicyID: policyID,
		Name:     name,
		Arn:      toArn(ArnPolicy, defaultPolicyPath+name),
		Document: doc,
	}

	query := `
	INSERT INTO policies (policy_id, name, arn, document)
	VALUES (:policy_id, :name, :arn, :document)
	`

	_, err := store.writeDB.NamedExec(query, &policy)
	if err != nil {
		slog.Error("create policy insert", "err", err)
		return nil, err
	}

	policy, err = GetPolicyByID(policyID)
	if err != nil {
		slog.Error("get policy by id", "err", err)
		return nil, &errs.ErrNotFound{Type: errs.ErrEntityTypePolicy, Resource: policyID}
	}

	return policy, nil
}

func CreatePolicyHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("PolicyName")
	if name == "" {
		errs.Handle(w, errors.New("missing query parameter: PolicyName"))
		return
	}

	document := r.URL.Query().Get("PolicyDocument")
	if document == "" {
		errs.Handle(w, errors.New("missing query parameter: PolicyDocument"))
		return
	}

	var doc PolicyDocument
	err := json.Unmarshal([]byte(document), &doc)
	if err != nil {
		errs.Handle(w, err)
		return
	}

	policy, err := CreatePolicy(name, &doc)
	if err != nil {
		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMServiceFailure())
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	policyPath := "/"

	response := &CreatePolicyResponse{
		CreatePolicyResult: &CreatePolicyResult{
			Policy: &PolicyResult{
				PolicyName:       policy.Name,
				DefaultVersionId: defaultPolicyVersion,
				PolicyId:         policy.PolicyID,
				Path:             policyPath,
				Arn:              policy.Arn,
				AttachmentCount:  policy.AttachmentCount,
				CreateDate:       policy.CreatedAt,
				UpdateDate:       policy.UpdatedAt,
			},
		},
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)

}

func AttachPolicy(arnType ArnType, policyArn, userName string) error {
	user, err := GetUserByName(userName)
	if err != nil {
		slog.Error("get user by name", "err", err)
		return err
	}

	policy, err := GetPolicyByArn(policyArn)
	if err != nil {
		slog.Error("get policy by arn", "err", err)
		return err
	}

	var tableName string
	var idFieldName string

	switch arnType {
	case ArnUser:
		tableName = "user_policies"
		idFieldName = "user_id"
	case ArnGroup:
		tableName = "group_policies"
		idFieldName = "group_id"
	default:
		return errors.New("unsupported arn type")
	}

	query_tmpl := `
	INSERT INTO %s (%s, policy_id)
	VALUES (:%s, :policy_id)
	`

	query := fmt.Sprintf(query_tmpl, tableName, idFieldName, idFieldName)

	_, err = store.writeDB.Exec(query, user.UserID, policy.PolicyID)
	if err != nil {
		slog.Error("attach policy insert", "err", err)
		return err
	}

	return nil
}
