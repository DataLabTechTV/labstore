package types

import "encoding/xml"

type UserResult struct {
	Path     string
	UserName string
	UserId   string
	Arn      string
}

type AccessKeyResult struct {
	XMLName         xml.Name `xml:"AccessKey"`
	UserName        string
	AccessKeyId     string
	Status          AccessKeyStatus
	SecretAccessKey string
}

type AccessKeyStatus string

const (
	AccessKeyActive   AccessKeyStatus = "Active"
	AccessKeyInactive AccessKeyStatus = "Inactive"
	AccessKeyExpired  AccessKeyStatus = "Expired"
)

type CreateUserResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ CreateUserResponse"`
	CreateUserResult *CreateUserResult
}

type CreateUserResult struct {
	User             *UserResult
	ResponseMetadata *ResponseMetadata
}

type CreateAccessKeyResponse struct {
	XMLName               xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ CreateAccessKeyResponse"`
	CreateAccessKeyResult *CreateAccessKeyResult
}

type CreateAccessKeyResult struct {
	AccessKey        *AccessKeyResult
	ResponseMetadata *ResponseMetadata
}

type GetUserResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ GetUserResponse"`
	GetUserResult    *GetUserResult
	ResponseMetadata *ResponseMetadata
}

type GetUserResult struct {
	User *UserResult
}

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

type DeleteUserResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DeleteUserResponse"`
	ResponseMetadata *ResponseMetadata
}

type DeleteAccessKeyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DeleteAccessKeyResponse"`
	ResponseMetadata *ResponseMetadata
}
