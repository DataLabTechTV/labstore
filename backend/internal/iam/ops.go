package iam

type IAMOp string

const (
	// ======================
	// USERS
	// ======================

	// --- Users: Create ---
	OpCreateUser       IAMOp = "CreateUser"
	OpCreateAccessKey  IAMOp = "CreateAccessKey"
	OpAttachUserPolicy IAMOp = "AttachUserPolicy"

	// --- Users: Read ---
	OpGetUser                  IAMOp = "GetUser"
	OpListAccessKeys           IAMOp = "ListAccessKeys"
	OpListAttachedUserPolicies IAMOp = "ListAttachedUserPolicies"

	// --- Users: Delete ---
	OpDeleteUser       IAMOp = "DeleteUser"
	OpDeleteAccessKey  IAMOp = "DeleteAccessKey"
	OpDetachUserPolicy IAMOp = "DetachUserPolicy"

	// ======================
	// GROUPS
	// ======================

	// --- Groups: Create ---
	OpCreateGroup       IAMOp = "CreateGroup"
	OpAddUserToGroup    IAMOp = "AddUserToGroup"
	OpAttachGroupPolicy IAMOp = "AttachGroupPolicy"

	// --- Groups: Read ---
	OpGetGroup                  IAMOp = "GetGroup"
	OpListAttachedGroupPolicies IAMOp = "ListAttachedGroupPolicies"

	// --- Groups: Delete ---
	OpDeleteGroup         IAMOp = "DeleteGroup"
	OpRemoveUserFromGroup IAMOp = "RemoveUserFromGroup"
	OpDetachGroupPolicy   IAMOp = "DetachGroupPolicy"

	// ======================
	// POLICIES
	// ======================

	// --- Policies: Create ---
	OpCreatePolicy IAMOp = "CreatePolicy"

	// --- Policies: Read ---
	OpGetPolicy IAMOp = "GetPolicy"

	// --- Policies: Delete ---
	OpDeletePolicy IAMOp = "DeletePolicy"
)
