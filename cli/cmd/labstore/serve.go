package main

import (
	"log/slog"

	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
	"github.com/IllumiKnowLabs/labstore/server/pkg/router"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "serve",
		Short: "Run server",
		Long:  "Run server for S3, IAM, and admin services",
		Run: func(cmd *cobra.Command, args []string) {
			adminPass := cmd.Flags().Lookup("admin-pass")
			if adminPass != nil && adminPass.Changed {
				slog.Warn("setting admin pass via the command line is insecure")
			}

			router.Start()
		},
	}

	cmd.Flags().String("storage-data-dir", config.DefaultStorageDataDir, "Storage root path for objects and metadata")
	cmd.Flags().String("storage-keys-dir", config.DefaultStorageDataDir, "Storage root path for encryption keys")

	cmd.Flags().String("admin-address-host", config.DefaultAdminAddressHost, "Listening host for admin server")
	cmd.Flags().String("admin-address-port", config.DefaultAdminAddressHost, "Listening port for admin server")
	cmd.Flags().String("admin-auth-access-key", config.DefaultAdminAuthAccessKey, "Administrator account access key")
	cmd.Flags().String("admin-auth-secret-key", config.DefaultAdminSecretKey, "Administrator account secret key")

	cmd.Flags().String("iam-address-host", config.DefaultIAMAddressHost, "Listening host for IAM server")
	cmd.Flags().String("iam-address-port", config.DefaultIAMAddressHost, "Listening port for IAM server")
	cmd.Flags().Int("iam-db-max-open-conns", config.DefaultIAMDBMaxOpenConns, "Maximum open reader connections for the IAM database")
	cmd.Flags().Int("iam-db-max-idle-conns", config.DefaultIAMDBMaxIdleConns, "Maximum idle reader connections for the IAM database")
	cmd.Flags().Int("iam-db-write-chan-cap", config.DefaultIAMWriteChanCap, "Buffered channel capacity for writing requests to the IAM database")
	cmd.Flags().Int("iam-db-timeout-ms", config.DefaultIAMDBTimeoutMs, "Connection timeout for the IAM database")
	cmd.Flags().Int("iam-db-read-cache-size-kib", config.DefaultIAMDBReadCacheSizeKiB, "Cache size of each individual reader connection for the IAM database")
	cmd.Flags().Int("iam-db-write-cache-size-kib", config.DefaultIAMDBWriteCacheSizeKiB, "Cache size of the writer connection for the IAM database")

	cmd.Flags().String("s3-address-host", config.DefaultS3AddressHost, "Listening host for S3-compatible server")
	cmd.Flags().Uint16("s3-address-port", config.DefaultS3AddressPort, "Listening port for S3-compatible server")
	cmd.Flags().Int("s3-paging-max-keys", config.DefaultS3PagingMaxKeys, "Hard limit for the maximum number of keys to return in paged requests")
	cmd.Flags().Int("s3-io-buffer-size", config.DefaultS3IOBufferSize, "Input/output buffer size in bytes")

	cmd.PersistentFlags().String("web-address-host", config.DefaultWebAddressHost, "Listening host for web ui server")
	cmd.PersistentFlags().Uint16("web-address-port", config.DefaultWebAddressPort, "Listening port for web ui server")

	return cmd
}
