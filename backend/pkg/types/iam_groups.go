package types

import "encoding/xml"

type GroupResult struct {
	Path      string
	GroupName string
	GroupId   string
	Arn       string
}

type CreateGroupResponse struct {
	XMLName           xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ CreateGroupResponse"`
	CreateGroupResult *CreateGroupResult
}

type CreateGroupResult struct {
	Group            *GroupResult
	ResponseMetadata *ResponseMetadata
}

type GetGroupResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ GetGroupResponse"`
	GetGroupResult   *GetGroupResult
	ResponseMetadata *ResponseMetadata
}

type GetGroupResult struct {
	Group       *GroupResult
	Users       *UserMembers
	IsTruncated bool
}

type UserMembers struct {
	Member []*UserResult `xml:"member"`
}

type DeleteGroupResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DeleteGroupResponse"`
	ResponseMetadata *ResponseMetadata
}
