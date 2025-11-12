# S3 Object Store API

## Planning

### Priority Levels

| Priority         | Description                                                                                                                                                                  |
| ---------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 🟥 P0 – Critical | Must be implemented first for the system to function or integrate—S3 clients, like `mc`, can run basic service, bucket, and object operations, including multi-part uploads. |
| 🟧 P1 – High     | Needed for MVP or user-visible value—IAM CRUD, including user/group management and policy management, and request authentication.                                            |
| 🟨 P2 – Medium   | Enhances usability or compliance—frequently used features, like versioning, encryption, locking, CORS, etc.                                                                  |
| 🟩 P3 – Low      | Nice-to-have, advanced, or admin-level—all other features.                                                                                                                   |

**Observation:** Entries for priorities 🟨 P2 – Medium and 🟩 P3 – Low are incomplete, underspecified, or underplanned, and will remain so until they become relevant.

### Development Status

| Status        | Description                                                                 |
| ------------- | --------------------------------------------------------------------------- |
| 🔴 Planned    | Work hasn’t started yet—the feature is on the roadmap.                      |
| 🟠 Developing | Actively developing—there is a dev branch or ongoing work.                  |
| 🟡 Testing    | Experimental or partial support—released for community testing (RC).        |
| 🟢 Released   | Fully implemented, tested, and deployed in production—no associated issues. |

## S3 REST API Endpoints

**Docs:** https://docs.aws.amazon.com/AmazonS3/latest/API/API_Operations_Amazon_Simple_Storage_Service.html

### SigV4

**Priority:** 🟥 P0 – Critical

| Spec            | Docs                                                                          |
| --------------- | ----------------------------------------------------------------------------- |
| Single Chunk    | https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-header-based-auth.html |
| Multiple Chunks | https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-streaming.html          |

### Service

**Priority:** 🟥 P0 – Critical

| S3 Action                                                                           | Method | Path | Description      | Status |
| ----------------------------------------------------------------------------------- | ------ | ---- | ---------------- | ------ |
| [ListBuckets](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListBuckets.html) | GET    | `/`  | List all buckets | 🟡     |

### Bucket

**Priority:** 🟥 P0 – Critical

| S3 Action                                                                               | Method | Path                    | Description                 | Status |
| --------------------------------------------------------------------------------------- | ------ | ----------------------- | --------------------------- | ------ |
| [CreateBucket](https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateBucket.html)   | PUT    | `/{bucket}`             | Create bucket               | 🟡     |
| [DeleteBucket](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucket.html)   | DELETE | `/{bucket}`             | Delete bucket               | 🟡     |
| [ListObjects](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjects.html)     | GET    | `/{bucket}`             | List objects in bucket      | 🔴     |
| [ListObjectsV2](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectsV2.html) | GET    | `/{bucket}?list-type=2` | List objects in bucket (V2) | 🔴     |
| [HeadBucket](https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadBucket.html)       | HEAD   | `/{bucket}`             | Check bucket existence      | 🔴     |

#### Configuration

**Priority:** 🟧 P1 – High

| S3 Action                                                                                                                                                                                                                                                                                     | Method         | Path                     | Description                   | Status |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------- | ------------------------ | ----------------------------- | ------ |
| [GetBucketAcl](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketAcl.html) / [PutBucketAcl](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketAcl.html)                                                                                                                 | GET/PUT        | `/{bucket}?acl`          | User/group permissions        | 🔴     |
| [GetBucketPolicy](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketPolicy.html) / [PutBucketPolicy](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketPolicy.html) / [DeleteBucketPolicy](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketPolicy.html) | GET/PUT/DELETE | `/{bucket}?policy`       | IAM-style JSON policy         | 🔴     |
| [GetBucketPolicyStatus](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketPolicyStatus.html)                                                                                                                                                                                       | GET            | `/{bucket}?policyStatus` | Check if the bucket is public | 🔴     |

**Priority:** 🟨 P2 – Medium

| S3 Action | Method         | Path                    | Description                                                           | Status |
| --------- | -------------- | ----------------------- | --------------------------------------------------------------------- | ------ |
|           | GET/PUT/DELETE | `/{bucket}?versioning`  | Configure object versioning                                           | 🔴     |
|           | GET            | `/{bucket}?versions`    | List all object versions                                              | 🔴     |
|           | GET/PUT/DELETE | `/{bucket}?encryption`  | Toggle encryption for new objects                                     | 🔴     |
|           | GET/PUT/DELETE | `/{bucket}?object-lock` | Configure object locks                                                | 🔴     |
|           | GET/PUT/DELETE | `/{bucket}?cors`        | CORS configurations to enable bucket operations from external domains | 🔴     |

**Priority:** 🟩 P3 – Low

| S3 Action | Method         | Path                     | Description                                                  | Status |
| --------- | -------------- | ------------------------ | ------------------------------------------------------------ | ------ |
|           | GET/PUT/DELETE | `/{bucket}?lifecycle`    | Time-based archival or deletion rules                        |        |
|           | GET/PUT/DELETE | `/{bucket}?replication`  | Automate bucket replication (role-based, filter-based, etc.) |        |
|           | GET/PUT/DELETE | `/{bucket}?website`      | Toggle bucket serving as a website                           |        |
|           | GET/PUT/DELETE | `/{bucket}?logging`      | Toggle request logging                                       |        |
|           | GET/PUT/DELETE | `/{bucket}?notification` |                                                              |        |
|           | GET/PUT/DELETE | `/{bucket}?metrics`      |                                                              |        |
|           | GET/PUT/DELETE | `/{bucket}?inventory`    |                                                              |        |
|           | GET/PUT/DELETE | `/{bucket}?accelerate`   |                                                              |        |
|           | GET/PUT/DELETE | `/{bucket}?analytics`    |                                                              |        |

### Object

**Priority:** 🟥 P0 – Critical

| S3 Action                                                                             | Method | Path                                    | Description        | Status |
| ------------------------------------------------------------------------------------- | ------ | --------------------------------------- | ------------------ | ------ |
| [PutObject](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html)       | PUT    | `/{bucket}/{key}`                       | Upload an object   | 🟡     |
| [GetObject](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html)       | GET    | `/{bucket}/{key}`                       | Download an object | 🟡     |
| [HeadObject](https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadObject.html)     | HEAD   | `/{bucket}/{key}`                       | Get metadata       | 🔴     |
| [DeleteObject](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html) | DELETE | `/{bucket}/{key}`                       | Delete an object   | 🟡     |
| [CopyObject](https://docs.aws.amazon.com/AmazonS3/latest/API/API_CopyObject.html)     | PUT    | `/{bucket}/{key}?x-amz-copy-source=...` | Copy object        | 🔴     |

#### Metadata and Tagging

**Priority:** 🟩 P3 – Low

| S3 Action | Method         | Path                      | Description | Status |
| --------- | -------------- | ------------------------- | ----------- | ------ |
|           | GET/PUT/DELETE | `/{bucket}/{key}?tagging` |             |        |
|           | GET            | `/{bucket}/{key}?torrent` |             |        |

#### Versioning and Retention

**Priority:** 🟩 P3 – Low

| S3 Action | Method         | Path                             | Description | Status |
| --------- | -------------- | -------------------------------- | ----------- | ------ |
|           | GET/PUT        | `/{bucket}/{key}?versionId={id}` |             |        |
|           | GET/PUT/DELETE | `/{bucket}/{key}?legal-hold`     |             |        |
|           | GET/PUT/DELETE | `/{bucket}/{key}?retention`      |             |        |

#### Restore or Select

**Priority:** 🟩 P3 – Low

| S3 Action | Method | Path                                   | Description | Status |
| --------- | ------ | -------------------------------------- | ----------- | ------ |
|           | POST   | `/{bucket}/{key}?restore`              |             |        |
|           | POST   | `/{bucket}/{key}?select&select-type=2` |             |        |

#### Presigned URLs

**Priority:** 🟩 P3 – Low

| S3 Action | Method | Path                                                                                                                | Description                                                       | Status |
| --------- | ------ | ------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- | ------ |
|           | GET    | `https://{bucket}.s3.amazonaws.com/{key}?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...&X-Amz-Signature=...` | Temporary access to an object without creating system credentials |        |

### Multipart Uploading

**Priority:** 🟧 P1 – High

Multipart uploading covers both the bucket and object scopes.

| S3 Action                                                                                                   | Method | Path                                           | Description                    | Status |
| ----------------------------------------------------------------------------------------------------------- | ------ | ---------------------------------------------- | ------------------------------ | ------ |
| [ListMultipartUploads](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListMultipartUploads.html)       | GET    | `/{bucket}?uploads`                            | List ongoing multipart uploads | 🔴     |
| [CreateMultipartUpload](https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateMultipartUpload.html)     | POST   | `/{bucket}/{key}?uploads`                      | Initiate upload                | 🔴     |
| [UploadPart](https://docs.aws.amazon.com/AmazonS3/latest/API/API_UploadPart.html)                           | PUT    | `/{bucket}/{key}?partNumber={n}&uploadId={id}` | Upload part                    | 🔴     |
| [ListParts](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListParts.html)                             | GET    | `/{bucket}?uploadId={id}`                      | List parts                     | 🔴     |
| [CompleteMultipartUpload](https://docs.aws.amazon.com/AmazonS3/latest/API/API_CompleteMultipartUpload.html) | POST   | `/{bucket}/{key}?uploadId={id}`                | Complete upload                | 🔴     |
| [AbortMultipartUpload](https://docs.aws.amazon.com/AmazonS3/latest/API/API_AbortMultipartUpload.html)       | DELETE | `/{bucket}/{key}?uploadId={id}`                | Abort upload                   | 🔴     |

**Priority:** 🟩 P3 – Low

| S3 Action                                                                                 | Method | Path                                           | Description                                                                 | Status |
| ----------------------------------------------------------------------------------------- | ------ | ---------------------------------------------- | --------------------------------------------------------------------------- | ------ |
| [UploadPartCopy](https://docs.aws.amazon.com/AmazonS3/latest/API/API_UploadPartCopy.html) | PUT    | `/{bucket}/{key}?partNumber={n}&uploadId={id}` | Upload part, extended with additional headers, to copy from existing bucket |        |

## IAM REST API Endpoints

**Docs:** https://docs.aws.amazon.com/IAM/latest/APIReference/API_Operations.html

This is implemented as a separate service, that can be called by the S3 REST endpoints as required. Here, a path prefix is something like `/division_abc/subdivision_xyz/engineering`.

**Priority:** 🟧 P1 – High

| IAM Action                                                                                            | Method | Path                          | Description                                                  | Status |
| ----------------------------------------------------------------------------------------------------- | ------ | ----------------------------- | ------------------------------------------------------------ | ------ |
| [CreateUser](https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateUser.html)                 | POST   | `/?Action=CreateUser`         | Create user (no password / can't login)                      | 🔴     |
| [CreateLoginProfile](https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateLoginProfile.html) | POST   | `/?Action=CreateLoginProfile` | Set password for user (can login)                            | 🔴     |
| [CreateAccessKey](https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateAccessKey.html)       | POST   | `/?Action=CreateAccessKey`    | Create user access key (API only auth / can't use for login) | 🔴     |
| [ListUsers](https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListUsers.html)                   | POST   | `/?Action=ListUsers`          | List users for a path prefix                                 | 🔴     |
| [GetUser](https://docs.aws.amazon.com/IAM/latest/APIReference/API_GetUser.html)                       | POST   | `/?Action=GetUser`            | Get user by name                                             | 🔴     |
| [DeleteUser](https://docs.aws.amazon.com/IAM/latest/APIReference/API_DeleteUser.html)                 | POST   | `/?Action=DeleteUser`         | Delete user by username                                      | 🔴     |
| [CreateGroup](https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateGroup.html)               | POST   | `/?Action=CreateGroup`        | Create group                                                 | 🔴     |
| [AddUserToGroup](https://docs.aws.amazon.com/IAM/latest/APIReference/API_AddUserToGroup.html)         | POST   | `/?Action=AddUserToGroup`     | Add user to group using names                                | 🔴     |
| [ListGroups](https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListGroups.html)                 | POST   | `/?Action=ListGroups`         | List groups for a path prefix                                | 🔴     |
| [ListGroupsForUser](https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListGroupsForUser.html)   | POST   | `/?Action=ListGroupsForUser`  | List groups for a given user                                 | 🔴     |
| [DeleteGroup](https://docs.aws.amazon.com/IAM/latest/APIReference/API_DeleteGroup.html)               | POST   | `/?Action=DeleteGroup`        | Delete group by name                                         | 🔴     |

**Priority:** 🟩 P3 – Low

| IAM Action                                                                                          | Method | Path                         | Description | Status |
| --------------------------------------------------------------------------------------------------- | ------ | ---------------------------- | ----------- | ------ |
| [AttachUserPolicy](https://docs.aws.amazon.com/IAM/latest/APIReference/API_AttachUserPolicy.html)   | POST   | `/?Action=AttachUserPolicy`  |             |        |
| [AttachGroupPolicy](https://docs.aws.amazon.com/IAM/latest/APIReference/API_AttachGroupPolicy.html) | POST   | `/?Action=AttachGroupPolicy` |             |        |
| ...                                                                                                 |        |                              |             |        |
