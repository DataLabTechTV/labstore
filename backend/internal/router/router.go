package router

import (
	"fmt"
	"os"

	"github.com/DataLabTechTV/labstore/backend/internal/bucket"
	"github.com/DataLabTechTV/labstore/backend/internal/config"
	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/DataLabTechTV/labstore/backend/internal/middleware"
	"github.com/DataLabTechTV/labstore/backend/internal/object"
	"github.com/DataLabTechTV/labstore/backend/internal/service"
	"github.com/DataLabTechTV/labstore/backend/pkg/iam"
	"github.com/DataLabTechTV/labstore/backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func Start() {
	os.MkdirAll(config.Env.StorageRoot, 0755)

	app := fiber.New()

	app.Use(middleware.AuthMiddleware())
	app.Use(middleware.IAMMiddleware())

	app.Put("/:bucket", middleware.WithIAM(iam.CreateBucket, bucket.PutBucketHandler))
	app.Put("/:bucket/:key", middleware.WithIAM(iam.PutObject, object.PutObjectHandler))

	app.Get("/", middleware.WithIAM(iam.ListAllMyBuckets, service.ListBucketsHandler))
	app.Get("/:bucket", middleware.WithIAM(iam.ListBucket, bucket.ListObjectsHandler))
	app.Get("/:bucket/:key", middleware.WithIAM(iam.GetObject, object.GetObjectHandler))

	app.Delete("/:bucket", middleware.WithIAM(iam.DeleteBucket, bucket.DeleteBucketHandler))
	app.Delete("/:bucket/:key", middleware.WithIAM(iam.DeleteObject, object.DeleteObjectHandler))

	app.Head("/:bucket", middleware.WithIAM(iam.ListBucket, bucket.HeadBucketHandler))
	app.Head("/:bucket/:key", middleware.WithIAM(iam.GetObject, object.HeadObjectHandler))

	app.Use(func(c *fiber.Ctx) error {
		core.HandleError(c, core.ErrorNotImplemented())
		return nil
	})

	port := fmt.Sprintf(":%d", config.Env.Port)
	logger.Log.Infoln("Starting minimal S3-compatible server on", port)
	logger.Log.Fatal(app.Listen(port))
}
