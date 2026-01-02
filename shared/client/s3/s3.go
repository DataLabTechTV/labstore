package s3

type S3Client struct {
	Host      string
	Port      uint16
	AccessKey string
	SecretKey string
}

func NewS3Client(host string, port uint16, accessKey, secretKey string) *S3Client {
	return &S3Client{
		Host:      host,
		Port:      port,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

func (client *S3Client) ListBuckets() ([]string, error) {
	// TODO: implement
	return []string{}, nil
}

func (client *S3Client) ListObjects(bucket, key string) ([]string, error) {
	// TODO: implement
	return []string{}, nil
}
