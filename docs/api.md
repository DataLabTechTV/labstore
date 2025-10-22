# REST API Endpoints

## Service

| Method | Path | Description                |
| ------ | ---- | -------------------------- |
| GET    | `/`  | List all buckets           |
| HEAD   | `/`  | Check service availability |

## Buckets

| Method | Path        | Description            |
| ------ | ----------- | ---------------------- |
| PUT    | `/{bucket}` | Create bucket          |
| DELETE | `/{bucket}` | Delete bucket          |
| GET    | `/{bucket}` | List objects in bucket |
| HEAD   | `/{bucket}` | Check bucket existence |

### Configuration

| Method         | Path                     | Description                                                           |
| -------------- | ------------------------ | --------------------------------------------------------------------- |
| GET/PUT/DELETE | `/{bucket}?acl`          | User/group permissions                                                |
| GET/PUT/DELETE | `/{bucket}?policy`       | IAM-style JSON policy                                                 |
| GET/PUT/DELETE | `/{bucket}?cors`         | CORS configurations to enable bucket operations from external domains |
| GET/PUT/DELETE | `/{bucket}?lifecycle`    | Time-based archival or deletion rules                                 |
| GET/PUT/DELETE | `/{bucket}?replication`  | Automate bucket replication (role-based, filter-based, etc.)          |
| GET/PUT/DELETE | `/{bucket}?website`      | Toggle bucket serving as a website                                    |
| GET/PUT/DELETE | `/{bucket}?logging`      | Toggle request logging                                                |
| GET/PUT/DELETE | `/{bucket}?encryption`   | Toggle encryption for new objects                                     |
| GET/PUT/DELETE | `/{bucket}?versioning`   | Configure object versioning                                           |
| GET/PUT/DELETE | `/{bucket}?notification` |                                                                       |
| GET/PUT/DELETE | `/{bucket}?metrics`      |                                                                       |
| GET/PUT/DELETE | `/{bucket}?inventory`    |                                                                       |
| GET/PUT/DELETE | `/{bucket}?accelerate`   |                                                                       |
| GET/PUT/DELETE | `/{bucket}?object-lock`  |                                                                       |
| GET/PUT/DELETE | `/{bucket}?analytics`    |                                                                       |

### Bucket Object Listing / Versions

| Method | Path                  | Description                    |
| ------ | --------------------- | ------------------------------ |
| GET    | `/{bucket}?list-type=2` | List objects (V2 API)          |
| GET    | `/{bucket}?versions`    | List all object versions       |
| GET    | `/{bucket}?uploads`     | List ongoing multipart uploads |

## Objects

| Method | Path                                  | Description        |
| ------ | ------------------------------------- | ------------------ |
| PUT    | `/{bucket}/{key}`                       | Upload an object   |
| GET    | `/{bucket}/{key}`                       | Download an object |
| HEAD   | `/{bucket}/{key}`                       | Get metadata       |
| DELETE | `/{bucket}/{key}`                       | Delete an object   |
| COPY   | `/{bucket}/{key}?x-amz-copy-source=...` | Copy object        |

### Metadata and Tagging

| Method         | Path                    | Description |
| -------------- | ----------------------- | ----------- |
| GET/PUT/DELETE | `/{bucket}/{key}?tagging` |             |
| GET/PUT        | `/{bucket}/{key}?acl`     |             |
| GET            | `/{bucket}/{key}?torrent` |             |

### Versioning and Retention

| Method         | Path                           | Description |
| -------------- | ------------------------------ | ----------- |
| GET/PUT        | `/{bucket}/{key}?versionId={id}` |             |
| GET/PUT/DELETE | `/{bucket}/{key}?legal-hold`     |             |
| GET/PUT/DELETE | `/{bucket}/{key}?retention`      |             |

### Restore or Select

| Method | Path                                 | Description |
| ------ | ------------------------------------ | ----------- |
| POST   | `/{bucket}/{key}?restore`              |             |
| POST   | `/{bucket}/{key}?select&select-type=2` |             |

## Multipart Uploads

| Method | Path                                           | Description     |
| ------ | ---------------------------------------------- | --------------- |
| POST   | `/{bucket}/{key}?uploads`                      | Initiate upload |
| PUT    | `/{bucket}/{key}?partNumber={n}&uploadId={id}` | Upload part     |
| GET    | `/{bucket}?uploadId={id}`                      | List parts      |
| POST   | `/{bucket}/{key}?uploadId={id}`                | Complete upload |
| DELETE | `/{bucket}/{key}?uploadId={id}`                | Abort upload    |

## Presigned URLs

| Method | Path                                                                                                                | Description |
| ------ | ------------------------------------------------------------------------------------------------------------------- | ----------- |
| GET    | `https://{bucket}.s3.amazonaws.com/{key}?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...&X-Amz-Signature=...` |             |

## IAM / Security

**Note:** Implemented as a separate service.

| Method | Path                        | Description |
| ------ | --------------------------- | ----------- |
| POST   | `/?Action=CreateUser`       |             |
| POST   | `/?Action=CreateAccessKey`  |             |
| POST   | `/?Action=AttachUserPolicy` |             |
| POST   | `/?Action=ListUsers`        |             |
| POST   | `/?Action=GetUser`          |             |
| POST   | `/?Action=DeleteUser`       |             |
