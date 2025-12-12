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
	OpGetUserPolicy            IAMOp = "GetUserPolicy"

	// --- Users: Update ---
	OpPutUserPolicy IAMOp = "PutUserPolicy"

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
	OpGetGroupPolicy            IAMOp = "GetGroupPolicy"

	// --- Groups: Update ---
	OpPutGroupPolicy IAMOp = "PutGroupPolicy"

	// --- Groups: Delete ---
	OpDeleteGroup         IAMOp = "DeleteGroup"
	OpRemoveUserFromGroup IAMOp = "RemoveUserFromGroup"
	OpDetachGroupPolicy   IAMOp = "DetachGroupPolicy"

	// ======================
	// POLICIES
	// ======================

	// --- Policies: Create ---
	OpCreatePolicy        IAMOp = "CreatePolicy"
	OpCreatePolicyVersion IAMOp = "CreatePolicyVersion"

	// --- Policies: Read ---
	OpGetPolicy        IAMOp = "GetPolicy"
	OpGetPolicyVersion IAMOp = "GetPolicyVersion"

	// --- Policies: Update ---
	// N/A

	// --- Policies: Delete ---
	OpDeletePolicy        IAMOp = "DeletePolicy"
	OpDeletePolicyVersion IAMOp = "DeletePolicyVersion"
)
