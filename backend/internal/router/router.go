package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/DataLabTechTV/labstore/backend/internal/bucket"
	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/middleware"
	"github.com/DataLabTechTV/labstore/backend/internal/object"
	"github.com/DataLabTechTV/labstore/backend/internal/service"
	"github.com/DataLabTechTV/labstore/backend/pkg/iam"
	"github.com/DataLabTechTV/labstore/backend/pkg/logger"
)

func Start() {
	os.MkdirAll(config.Env.StorageRoot, 0755)

	mux := http.NewServeMux()

	mux.Handle("PUT /{bucket}", middleware.WithIAM(iam.CreateBucket, http.HandlerFunc(bucket.PutBucketHandler)))
	mux.Handle("PUT /{bucket}/{key...}", middleware.WithIAM(iam.PutObject, http.HandlerFunc(object.PutObjectHandler)))

	mux.Handle("GET /", middleware.WithIAM(iam.ListAllMyBuckets, http.HandlerFunc(service.ListBucketsHandler)))
	mux.Handle("GET /{bucket}", middleware.WithIAM(iam.ListBucket, http.HandlerFunc(bucket.ListObjectsHandler)))
	mux.Handle("GET /{bucket}/{key...}", middleware.WithIAM(iam.GetObject, http.HandlerFunc(object.GetObjectHandler)))

	mux.Handle("DELETE /{bucket}", middleware.WithIAM(iam.DeleteBucket, http.HandlerFunc(bucket.DeleteBucketHandler)))
	mux.Handle("DELETE /{bucket}/{key...}", middleware.WithIAM(iam.DeleteObject, http.HandlerFunc(object.DeleteObjectHandler)))

	mux.Handle("HEAD /{bucket}", middleware.WithIAM(iam.ListBucket, http.HandlerFunc(bucket.HeadBucketHandler)))
	mux.Handle("HEAD /{bucket}/{key...}", middleware.WithIAM(iam.GetObject, http.HandlerFunc(object.HeadObjectHandler)))

	addr := fmt.Sprintf("%s:%d", config.Env.Host, config.Env.Port)
	logger.Log.Infoln("Starting minimal S3-compatible server on", addr)

	server := http.Server{
		Addr: addr,
		Handler: chain(
			mux,
			middleware.CompressionMiddleware,
			middleware.AuthMiddleware,
			middleware.IAMMiddleware,
			middleware.NormalizeMiddleware,
		),
	}

	logger.Log.Fatal(server.ListenAndServe())
}

func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
