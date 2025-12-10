package iam

type IAMOp string

// TODO: remaining ops for CRUD
const (
	// Users
	OpCreateUser      IAMOp = "CreateUser"
	OpCreateAccessKey IAMOp = "CreateAccessKey"

	// Groups
	OpCreateGroup IAMOp = "CreateGroup"

	// Policies
	OpCreatePolicy      IAMOp = "CreatePolicy"
	OpAttachUserPolicy  IAMOp = "AttachUserPolicy"
	OpAttachGroupPolicy IAMOp = "AttachGroupPolicy"
)
