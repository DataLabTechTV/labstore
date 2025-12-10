package router

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/bucket"
	"github.com/IllumiKnowLabs/labstore/backend/internal/iam"
	"github.com/IllumiKnowLabs/labstore/backend/internal/middleware"
	"github.com/IllumiKnowLabs/labstore/backend/internal/object"
	"github.com/IllumiKnowLabs/labstore/backend/internal/service"
)

func NewS3ServerDescriptor(host string, port uint16) *ServerDescriptor {
	slog.Info(
		"s3-compatible api server",
		"host", host,
		"port", port,
	)

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
	router.Handle("GET /", middleware.WithIAM(iam.ListAllMyBuckets, http.HandlerFunc(service.ListBucketsHandler)))

	// Bucket
	router.Handle("HEAD /{bucket}", middleware.WithIAM(iam.ListBucket, http.HandlerFunc(bucket.HeadBucketHandler)))
	router.Handle("GET /{bucket}", middleware.WithIAM(iam.ListBucket, http.HandlerFunc(bucket.ListObjectsHandler)))
	router.Handle("PUT /{bucket}", middleware.WithIAM(iam.CreateBucket, http.HandlerFunc(bucket.PutBucketHandler)))
	router.Handle("DELETE /{bucket}", middleware.WithIAM(iam.DeleteBucket, http.HandlerFunc(bucket.DeleteBucketHandler)))

	// Object
	router.Handle("HEAD /{bucket}/{key...}", middleware.WithIAM(iam.GetObject, http.HandlerFunc(object.HeadObjectHandler)))
	router.Handle("GET /{bucket}/{key...}", middleware.WithIAM(iam.GetObject, http.HandlerFunc(object.GetObjectHandler)))
	router.Handle("PUT /{bucket}/{key...}", middleware.WithIAM(iam.PutObject, http.HandlerFunc(object.PutObjectHandler)))
	router.Handle("DELETE /{bucket}/{key...}", middleware.WithIAM(iam.DeleteObject, http.HandlerFunc(object.DeleteObjectHandler)))
	router.Handle("POST /{bucket}", middleware.WithIAM(iam.DeleteBucket, http.HandlerFunc(object.DeleteObjectsHandler)))
}
