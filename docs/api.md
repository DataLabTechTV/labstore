# S3 Object Store API

## Planning

### Priority Levels

| Priority | Description |
| ---------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 🟥 P0 – Critical | Must be implemented first for the system to function or integrate—S3 clients, like `mc`, can run basic service, bucket, and object operations, including multi-part uploads. |
| 🟧 P1 – High | Needed for MVP or user-visible value—IAM CRUD, including user/group management and policy management, and request authentication. |
| 🟨 P2 – Medium | Enhances usability or compliance—frequently used features, like versioning, encryption, locking, CORS, etc. |
| 🟩 P3 – Low | Nice-to-have, advanced, or admin-level—all other features. |

**Observation:** Entries for priorities 🟨 P2 – Medium and 🟩 P3 – Low are incomplete, underspecified, or underplanned, and will remain so until they become relevant.

### Development Status

| Status | Description |
| ------------- | --------------------------------------------------------------------------- |
| 🔴 Planned | Work hasn’t started yet—the feature is on the roadmap. |
| 🟠 Developing | Actively developing—there is a dev branch or ongoing work. |
| 🟡 Testing | Experimental or partial support—released for community testing (RC). |
| 🟢 Released | Fully implemented, tested, and deployed in production—no associated issues. |

## S3 REST API Endpoints

**Docs:** https://docs.aws.amazon.com/AmazonS3/latest/API/API_Operations_Amazon_Simple_Storage_Service.html

### SigV4

**Priority:** 🟥 P0 – Critical

| Spec | Status |
| ---- | ------ |
| [SigV4 - Single Chunk](https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-header-based-auth.html) | 🟢 |
| [SigV4 - Multiple Chunks](https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-streaming.html) | 🟢 |

### Service

**Priority:** 🟥 P0 – Critical

| S3 Action | Method | Path | Description | Status |
| ----------------------------------------------------------------------------------- | ------ | ---- | ---------------- | ------ |
| [ListBuckets](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListBuckets.html) | GET | `/` | List all my buckets | 🟡 |

### Bucket

**Priority:** 🟥 P0 – Critical

| S3 Action | Method | Path | Description | Status |
| --------------------------------------------------------------------------------------- | ------ | ----------------------- | --------------------------- | ------ |
| [HeadBucket](https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadBucket.html) | HEAD | `/{bucket}` | Check bucket existence | 🟡 |
| [ListObjects](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjects.html) | GET | `/{bucket}` | List objects in bucket | 🟡 |
| [ListObjectsV2](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectsV2.html) | GET | `/{bucket}?list-type=2` | List objects in bucket (V2) | 🟡 |
| [CreateBucket](https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateBucket.html) | PUT | `/{bucket}` | Create bucket | 🟡 |
| [DeleteBucket](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucket.html) | DELETE | `/{bucket}` | Delete bucket | 🟡 |

#### Configuration

**Priority:** 🟧 P1 – High

| S3 Action | Method | Path | Description | Status |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------- | ------------------------ | ----------------------------- | ------ |
| [GetBucketAcl](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketAcl.html) / [PutBucketAcl](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketAcl.html) | GET/PUT | `/{bucket}?acl` | User/group permissions | 🔴 |
| [GetBucketPolicy](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketPolicy.html) / [PutBucketPolicy](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketPolicy.html) / [DeleteBucketPolicy](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketPolicy.html) | GET/PUT/DELETE | `/{bucket}?policy` | IAM-style JSON policy | 🔴 |
| [GetBucketPolicyStatus](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketPolicyStatus.html) | GET | `/{bucket}?policyStatus` | Check if the bucket is public | 🔴 |

**Priority:** 🟨 P2 – Medium

| S3 Action | Method | Path | Description | Status |
| --------- | -------------- | ----------------------- | --------------------------------------------------------------------- | ------ |
| [GetBucketVersioning](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketVersioning.html) / [PutBucketVersioning](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketVersioning.html) /  | GET/PUT | `/{bucket}?versioning` | Configure object versioning | 🔴 |
| [ListObjectVersions](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectVersions.html) | GET | `/{bucket}?versions` | List all object versions | 🔴 |
| [GetBucketEncryption](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketEncryption.html) / [PutBucketEncryption](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketEncryption.html) / [DeleteBucketEncryption](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketEncryption.html) | GET/PUT/DELETE | `/{bucket}?encryption` | Toggle encryption for new objects | 🔴 |
| [GetObjectLockConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObjectLockConfiguration.html) / [PutObjectLockConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObjectLockConfiguration.html) | GET/PUT | `/{bucket}?object-lock` | Configure object locks | 🔴 |
| [GetBucketCors](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketCors.html) / [PutBucketCors](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketCors.html) / [DeleteBucketCors](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketCors.html) | GET/PUT/DELETE | `/{bucket}?cors` | CORS configurations to enable bucket operations from external domains | 🔴 |

**Priority:** 🟩 P3 – Low

| S3 Action | Method | Path | Description | Status |
| --------- | -------------- | ------------------------ | ------------------------------------------------------------ | ------ |
| [GetBucketLifecycle](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketLifecycle.html) / [PutBucketLifecycle](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketLifecycle.html) | GET/PUT | `/{bucket}?lifecycle` | Time-based archival or deletion rules | 🔴 |
| [GetBucketReplication](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketReplication.html) / [PutBucketReplication](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketReplication.html) / [DeleteBucketReplication](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketReplication.html) | GET/PUT/DELETE | `/{bucket}?replication` | Automate bucket replication (role-based, filter-based, etc.) | 🔴 |
| [GetWebsite](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetWebsite.html) / [PutWebsite](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutWebsite.html) / [DeleteWebsite](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteWebsite.html) | GET/PUT/DELETE | `/{bucket}?website` | Manage bucket serving as a website | 🔴 |
| [GetBucketLogging](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketLogging.html) / [PutBucketLogging](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketLogging.html) | GET/PUT | `/{bucket}?logging` | Toggle request logging | 🔴 |
| [GetBucketNotificationConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketNotificationConfiguration.html) / [PutBucketNotificationConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketNotificationConfiguration.html) | GET/PUT | `/{bucket}?notification` | Manage topic, queue, and cloud functions notifications | 🔴 |
| [GetBucketMetricsConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketMetricsConfiguration.html) / [PutBucketMetricsConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketMetricsConfiguration.html) / [DeleteBucketMetricsConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketMetricsConfiguration.html) | GET/PUT/DELETE | `/{bucket}?metrics` | Manage metrics configuration to be used with Amazon CloudWatch | 🔴 |
| [GetBucketInventoryConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketInventoryConfiguration.html) / [ListBucketInventoryConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListBucketInventoryConfiguration.html) / [PutBucketInventoryConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketInventoryConfiguration.html) / [DeleteBucketInventoryConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketInventoryConfiguration.html) | GET/LIST/PUT/DELETE | `/{bucket}?inventory` | Manage inventory report configurations | 🔴 |
| [GetBucketAccelerateConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketAccelerateConfiguration.html) / [PutBucketAccelerateConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketAccelerateConfiguration.html) | GET/PUT | `/{bucket}?accelerate` | Toggles bucket data acceleration for faster transfer | 🔴 |
| [GetBucketAnalyticsConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketAnalyticsConfiguration.html) / [ListBucketAnalyticsConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListBucketAnalyticsConfiguration.html) / [PutBucketAnalyticsConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketAnalyticsConfiguration.html) / [DeleteBucketAnalyticsConfiguration](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketAnalyticsConfiguration.html) | GET/LIST/PUT/DELETE | `/{bucket}?analytics` | Manage analytics configuration that help monitor access patterns and storage usage | 🔴 |

### Object

**Priority:** 🟥 P0 – Critical

| S3 Action | Method | Path | Description | Status |
| ------------------------------------------------------------------------------------- | ------ | --------------------------------------- | ------------------ | ------ |
| [HeadObject](https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadObject.html) | HEAD | `/{bucket}/{key}` | Get metadata | 🟡 |
| [GetObject](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html) | GET | `/{bucket}/{key}` | Download an object | 🟡 |
| [PutObject](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html) | PUT | `/{bucket}/{key}` | Upload an object | 🟡 |
| [DeleteObject](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html) | DELETE | `/{bucket}/{key}` | Delete an object | 🟡 |
| [DeleteObjects](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObjects.html) | POST | `/{bucket}?delete` | Delete multiple objects (batch delete) | 🟡 |

**Priority:** 🟧 P1 – High

| S3 Action | Method | Path | Description | Status |
| ------------------------------------------------------------------------------------- | ------ | --------------------------------------- | ------------------ | ------ |
| [CopyObject](https://docs.aws.amazon.com/AmazonS3/latest/API/API_CopyObject.html) | PUT | `/{bucket}/{key}?x-amz-copy-source=...` | Copy object | 🔴 |

#### Metadata and Tagging

**Priority:** 🟩 P3 – Low

| S3 Action | Method | Path | Description | Status |
| --------- | -------------- | ------------------------- | ----------- | ------ |
| [GetObjectTagging](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObjectTagging.html) / [PutObjectTagging](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObjectTagging.html) / [DeleteObjectTagging](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObjectTagging.html) | GET/PUT/DELETE | `/{bucket}/{key}?tagging` | Manage object key-value style tags | 🔴 |
| [GetObjectTorrent](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObjectTorrent.html) | GET | `/{bucket}/{key}?torrent` | BitTorrent files for a bucket key | 🔴 |

#### Versioning and Retention

**Priority:** 🟩 P3 – Low

| S3 Action | Method | Path | Description | Status |
| --------- | -------------- | -------------------------------- | ----------- | ------ |
| [GetObject](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html) | GET | `/{bucket}/{key}?versionId={id}` | Download a specific version of an object | 🔴 |
| [GetObjectLegalHold](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObjectLegalHold.html) / [PutObjectLegalHold](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObjectLegalHold.html) | GET/PUT | `/{bucket}/{key}?legal-hold` | Toggle legal hold status for an object | 🔴 |
| [GetObjectRetention](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObjectRetention.html) / [PutObjectRetention](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObjectRetention.html) | GET/PUT | `/{bucket}/{key}?retention` | Gets or sets object retention up to a given date | 🔴 |

#### Restore or Select

**Priority:** 🟩 P3 – Low

| S3 Action | Method | Path | Description | Status |
| --------- | ------ | -------------------------------------- | ----------- | ------ |
| [RestoreObject](https://docs.aws.amazon.com/AmazonS3/latest/API/API_RestoreObject.html) | POST | `/{bucket}/{key}?restore` | Restores an archived object | 🔴 |
| [SelectObjectContent](https://docs.aws.amazon.com/AmazonS3/latest/API/API_SelectObjectContent.html) | POST | `/{bucket}/{key}?select&select-type=2` | Directly query JSON, CSV or Parquet objects using SQL | 🔴 |

#### Presigned URLs

**Priority:** 🟩 P3 – Low

| Docs | Method | Path | Description | Status |
| --------- | ------ | ------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- | ------ |
| [Presigned URLs](https://docs.aws.amazon.com/AmazonS3/latest/userguide/using-presigned-url.html) | GET | `/{bucket}/{key}?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...&X-Amz-Signature=...` | Temporary access to an object without creating system credentials | 🔴 |

### Multipart Uploading

**Priority:** 🟧 P1 – High

Multipart uploading covers both the bucket and object scopes.

| S3 Action | Method | Path | Description | Status |
| ----------------------------------------------------------------------------------------------------------- | ------ | ---------------------------------------------- | ------------------------------ | ------ |
| [ListMultipartUploads](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListMultipartUploads.html) | GET | `/{bucket}?uploads` | List ongoing multipart uploads | 🔴 |
| [CreateMultipartUpload](https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateMultipartUpload.html) | POST | `/{bucket}/{key}?uploads` | Initiate upload | 🔴 |
| [UploadPart](https://docs.aws.amazon.com/AmazonS3/latest/API/API_UploadPart.html) | PUT | `/{bucket}/{key}?partNumber={n}&uploadId={id}` | Upload part | 🔴 |
| [ListParts](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListParts.html) | GET | `/{bucket}?uploadId={id}` | List parts | 🔴 |
| [CompleteMultipartUpload](https://docs.aws.amazon.com/AmazonS3/latest/API/API_CompleteMultipartUpload.html) | POST | `/{bucket}/{key}?uploadId={id}` | Complete upload | 🔴 |
| [AbortMultipartUpload](https://docs.aws.amazon.com/AmazonS3/latest/API/API_AbortMultipartUpload.html) | DELETE | `/{bucket}/{key}?uploadId={id}` | Abort upload | 🔴 |

**Priority:** 🟩 P3 – Low

| S3 Action | Method | Path | Description | Status |
| ----------------------------------------------------------------------------------------- | ------ | ---------------------------------------------- | --------------------------------------------------------------------------- | ------ |
| [UploadPartCopy](https://docs.aws.amazon.com/AmazonS3/latest/API/API_UploadPartCopy.html) | PUT | `/{bucket}/{key}?partNumber={n}&uploadId={id}` | Upload part, extended with additional headers, to copy from existing bucket | |

## IAM REST API Endpoints

**Docs:** https://docs.aws.amazon.com/IAM/latest/APIReference/API_Operations.html

This is implemented as a separate service, that can be called by the S3 REST endpoints as required. Here, a path prefix is something like `/division_abc/subdivision_xyz/engineering`.

**Priority:** 🟧 P1 – High

| IAM Action | Method | Path | Description | Status |
| ----------------------------------------------------------------------------------------------------- | ------ | ----------------------------- | ------------------------------------------------------------ | ------ |
| [CreateUser](https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateUser.html) | POST | `/?Action=CreateUser` | Create user (no password / can't login) | 🟡 |
| [CreateAccessKey](https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateAccessKey.html) | POST | `/?Action=CreateAccessKey` | Create user access key (API only auth / can't use for login) | 🟡 |
| [AttachUserPolicy](https://docs.aws.amazon.com/IAM/latest/APIReference/API_AttachUserPolicy.html) | POST | `/?Action=AttachUserPolicy` | Attach a managed policy to a user | 🟡 |
| [GetUser](https://docs.aws.amazon.com/IAM/latest/APIReference/API_GetUser.html) | POST | `/?Action=GetUser` | Get user by user name | 🟡 |
| [ListAccessKeys](https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListAccessKeys.html) | POST | `/?Action=ListAccessKeys` | List access keys for a given user by user name | 🟡 |
| [ListAttachedUserPolicies](https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListAttachedUserPolicies.html) | POST | `/?Action=ListAttachedUserPolicies` | List managed policies for a given user by user name | 🟡 |
| [DeleteUser](https://docs.aws.amazon.com/IAM/latest/APIReference/API_DeleteUser.html) | POST | `/?Action=DeleteUser` | Delete user by user name | 🟡 |
| [DeleteAccessKey](https://docs.aws.amazon.com/IAM/latest/APIReference/API_DeleteAccessKey.html) | POST | `/?Action=DeleteAccessKey` | Delete access key by user name and access key ID | 🟡 |
| [DetachUserPolicy](https://docs.aws.amazon.com/IAM/latest/APIReference/API_DetachUserPolicy.html) | POST | `/?Action=DetachUserPolicy` | Detach managed policy from user by user name and policy ARN | 🟡 |
| [CreateGroup](https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateGroup.html) | POST | `/?Action=CreateGroup` | Create group | 🟡 |
| [AddUserToGroup](https://docs.aws.amazon.com/IAM/latest/APIReference/API_AddUserToGroup.html) | POST | `/?Action=AddUserToGroup` | Add user to group by user name and group name | 🟡 |
| [AttachGroupPolicy](https://docs.aws.amazon.com/IAM/latest/APIReference/API_AttachGroupPolicy.html) | POST | `/?Action=AttachGroupPolicy` | Attach a managed policy to a group by group name and policy ARN | 🟡 |
| [GetGroup](https://docs.aws.amazon.com/IAM/latest/APIReference/API_GetGroup.html) | POST | `/?Action=GetGroup` | Get group by group name | 🟡 |
| [ListAttachedGroupPolicies](https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListAttachedGroupPolicies.html) | POST | `/?Action=ListAttachedGroupPolicies` | List managed policies attached to a group by group name | 🟡 |
| [DeleteGroup](https://docs.aws.amazon.com/IAM/latest/APIReference/API_DeleteGroup.html) | POST | `/?Action=DeleteGroup` | Delete group by group name | 🟡 |
| [RemoveUserFromGroup](https://docs.aws.amazon.com/IAM/latest/APIReference/API_RemoveUserFromGroup.html) | POST | `/?Action=RemoveUserFromGroup` | Remove user from group by user name and group name | 🟡 |
| [DetachUserGroup](https://docs.aws.amazon.com/IAM/latest/APIReference/API_DetachUserGroup.html) | POST | `/?Action=DetachUserGroup` | Detach managed policy from group by group name and policy ARN | 🟡 |

**Priority:** 🟨 P2 – Medium

| IAM Action | Method | Path | Description | Status |
| --------------------------------------------------------------------------------------------------- | ------ | ---------------------------- | ----------- | ------ |
| [ListUsers](https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListUsers.html) | POST | `/?Action=ListUsers` | List users for a path prefix | 🔴 |
| [ListGroups](https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListGroups.html) | POST | `/?Action=ListGroups` | List groups for a path prefix | 🔴 |
| [ListGroupsForUser](https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListGroupsForUser.html) | POST | `/?Action=ListGroupsForUser` | List groups for a given user | 🔴 |

**Priority:** 🟩 P3 – Low

| IAM Action | Method | Path | Description | Status |
| --------------------------------------------------------------------------------------------------- | ------ | ---------------------------- | ----------- | ------ |
| [CreateLoginProfile](https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateLoginProfile.html) | POST | `/?Action=CreateLoginProfile` | Set password for user (can login) | 🔴 |
