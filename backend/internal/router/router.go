package router

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/IllumiKnowLabs/labstore/backend/internal/bucket"
	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/iam"
	"github.com/IllumiKnowLabs/labstore/backend/internal/middleware"
	"github.com/IllumiKnowLabs/labstore/backend/internal/object"
	"github.com/IllumiKnowLabs/labstore/backend/internal/service"
)

func Start() {
	if err := ensureDirectories(); err != nil {
		slog.Error("could not create directory structure", "err", err)
		os.Exit(1)
	}

	router := http.NewServeMux()
	loadRoutes(router)

	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)

	mw := middleware.Stack(
		middleware.LoggingMiddleware,
		middleware.CompressionMiddleware,
		middleware.AuthMiddleware,
		middleware.NormalizationMiddleware,
	)

	slog.Info(
		"starting S3-compatible object store server",
		"host", config.Server.Host,
		"port", config.Server.Port,
	)

	server := http.Server{
		Addr:    addr,
		Handler: mw(router),
	}

	fmt.Printf("🌐 Backend listening on http://%s\n", addr)

	log.Fatal(server.ListenAndServe())
}

func ensureDirectories() error {
	slog.Debug("ensuring directories")

	if err := os.MkdirAll(config.Server.Storage.Path, 0755); err != nil {
		return err
	}

	return nil
}

func loadRoutes(router *http.ServeMux) {
	slog.Debug("loading routes")

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
