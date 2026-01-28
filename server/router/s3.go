package router

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/server/buckets"
	"github.com/IllumiKnowLabs/labstore/server/iam"
	"github.com/IllumiKnowLabs/labstore/server/middleware"
	"github.com/IllumiKnowLabs/labstore/server/objects"
	"github.com/IllumiKnowLabs/labstore/server/service"
)

func NewS3ServerDescriptor(host string, port uint16) *ServerDescriptor {
	slog.Info("s3-compatible api server", "host", host, "port", port)

	router := http.NewServeMux()
	loadS3Routes(router)

	addr := fmt.Sprintf("%s:%d", host, port)

	mw := middleware.Stack(
		middleware.LoggingMiddleware,
		middleware.CompressionMiddleware,
		middleware.AuthMiddleware,
		middleware.NormalizationMiddleware,
	)

	server := http.Server{
		Addr:    addr,
		Handler: mw(router),
	}

	serverDescriptor := &ServerDescriptor{
		Name:   "S3-Compatible API",
		Server: &server,
	}

	return serverDescriptor
}

func loadS3Routes(router *http.ServeMux) {
	slog.Debug("loading s3-compatible routes")

	// Service
	router.Handle("GET /", middleware.WithIAM(iam.S3ListAllMyBuckets, http.HandlerFunc(service.ListBucketsHandler)))

	// Bucket
	router.Handle("HEAD /{bucket}", middleware.WithIAM(iam.S3ListBucket, http.HandlerFunc(buckets.HeadBucketHandler)))
	router.Handle("GET /{bucket}", middleware.WithIAM(iam.S3ListBucket, http.HandlerFunc(buckets.ListObjectsHandler)))
	router.Handle("PUT /{bucket}", middleware.WithIAM(iam.S3CreateBucket, http.HandlerFunc(buckets.PutBucketHandler)))
	router.Handle("DELETE /{bucket}", middleware.WithIAM(iam.S3DeleteBucket, http.HandlerFunc(buckets.DeleteBucketHandler)))

	// Object
	router.Handle("HEAD /{bucket}/{key...}", middleware.WithIAM(iam.S3GetObject, http.HandlerFunc(objects.HeadObjectHandler)))
	router.Handle("GET /{bucket}/{key...}", middleware.WithIAM(iam.S3GetObject, http.HandlerFunc(objects.GetObjectHandler)))
	router.Handle("PUT /{bucket}/{key...}", middleware.WithIAM(iam.S3PutObject, http.HandlerFunc(objects.PutObjectHandler)))
	router.Handle("DELETE /{bucket}/{key...}", middleware.WithIAM(iam.S3DeleteObject, http.HandlerFunc(objects.DeleteObjectHandler)))
	router.Handle("POST /{bucket}", middleware.WithIAM(iam.S3DeleteBucket, http.HandlerFunc(objects.DeleteObjectsHandler)))
}
