package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/IllumiKnowLabs/labstore/backend/internal/constants"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
	"github.com/IllumiKnowLabs/labstore/backend/internal/security"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configBasename = "labstore"

	DefaultAdminServerHost    = "0.0.0.0"
	DefaultAdminServerPort    = 6787
	DefaultAdminAuthAccessKey = "admin"
	DefaultAdminSecretKey     = DefaultAdminAuthAccessKey

	DefaultIAMServerHost = "0.0.0.0"
	DefaultIAMServerPort = 6788

	DefaultS3ServerHost     = "0.0.0.0"
	DefaultS3ServerPort     = 6789
	DefaultS3StoragePath    = "./data"
	DefaultS3PerfBufferSize = 256 * helper.KiB
)

type AppConfig struct {
	Backend *BackendConfig `mapstructure:"backend"`
}

type BackendConfig struct {
	Admin *AdminConfig `mapstructure:"admin"`
	IAM   *IAMConfig   `mapstructure:"iam"`
	S3    *S3Config    `mapstructure:"s3"`
}

type AdminConfig struct {
	Server *ServerConfig `mapstructure:"server"`
	Auth   *AuthConfig   `mapstructure:"auth"`
}

type IAMConfig struct {
	Server *ServerConfig `mapstructure:"server"`
}

type S3Config struct {
	Server  *ServerConfig  `mapstructure:"server"`
	Storage *StorageConfig `mapstructure:"storage"`
	Perf    *PerfConfig    `mapstructure:"perf"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

type AuthConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type StorageConfig struct {
	Path string `mapstructure:"path"`
}

type PerfConfig struct {
	BufferSize int `mapstructure:"buffer_size"`
}

var Admin *AdminConfig
var IAM *IAMConfig
var S3 *S3Config

func (config *AdminConfig) Debug() {
	slog.Debug("config set", "name", "backend.admin.server.host", "value", config.Server.Host)
	slog.Debug("config set", "name", "backend.admin.server.port", "value", config.Server.Port)
	slog.Debug("config set", "name", "backend.admin.access_key", "value", config.Auth.AccessKey)

	var adminSecretKeyDisplay string
	if len(config.Auth.SecretKey) > 0 {
		adminSecretKeyDisplay = security.Redacted
	} else {
		adminSecretKeyDisplay = constants.Empty
	}
	slog.Debug("config set", "name", "backend.admin.secret_key", "value", adminSecretKeyDisplay)
}

func (config *IAMConfig) Debug() {
	slog.Debug("config set", "name", "backend.iam.server.host", "value", config.Server.Host)
	slog.Debug("config set", "name", "backend.iam.server.port", "value", config.Server.Port)
}

func (config *S3Config) Debug() {
	slog.Debug("config set", "name", "backend.s3.server.host", "value", config.Server.Host)
	slog.Debug("config set", "name", "backend.s3.server.port", "value", config.Server.Port)
	slog.Debug("config set", "name", "backend.s3.storage.path", "value", config.Storage.Path)
	slog.Debug("config set", "name", "backend.s3.perf.buffer_size", "value", config.Perf.BufferSize)
}

func Load(rootCmd *cobra.Command) {
	slog.Debug("loading config")
	setLookupPaths()
	setDefaults()
	readConfig()
	setOverrides(rootCmd)
	parseConfig()
}

func setLookupPaths() {
	slog.Debug("setting config lookup paths")

	viper.SetConfigName(configBasename)

	viper.AddConfigPath(".")

	systemConfigPath := filepath.Join("/etc", configBasename)
	viper.AddConfigPath(systemConfigPath)

	if configDir, err := os.UserConfigDir(); err != nil {
		slog.Warn("user config dir skipped", "err", err)
	} else {
		userConfigPath := filepath.Join(configDir, configBasename)
		viper.AddConfigPath(userConfigPath)
	}
}

func setDefaults() {
	slog.Debug("setting config defaults")

	// admin
	viper.SetDefault("backend.admin.server.host", DefaultAdminServerHost)
	viper.SetDefault("backend.admin.server.port", DefaultAdminServerPort)
	viper.SetDefault("backend.admin.auth.access_key", DefaultAdminAuthAccessKey)

	defaultAdminAuthSecretKey, err := security.GeneratePassword(32)
	if err != nil {
		slog.Error("admin password generation", "err", err)
		defaultAdminAuthSecretKey = DefaultAdminSecretKey
	}

	viper.SetDefault("backend.admin.auth.secret_key", defaultAdminAuthSecretKey)
	fmt.Printf("🔑 Default admin secret key: %s\n", defaultAdminAuthSecretKey)

	// iam
	viper.SetDefault("backend.iam.server.host", DefaultIAMServerHost)
	viper.SetDefault("backend.iam.server.port", DefaultIAMServerPort)

	// s3
	viper.SetDefault("backend.s3.server.host", DefaultS3ServerHost)
	viper.SetDefault("backend.s3.server.port", DefaultS3ServerPort)
	viper.SetDefault("backend.s3.storage.path", DefaultS3StoragePath)
	viper.SetDefault("backend.s3.perf.buffer_size", DefaultS3PerfBufferSize)
}

func setOverrides(rootCmd *cobra.Command) {
	slog.Debug("setting config env and cli overrides")

	viper.SetEnvPrefix(configBasename)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	serverCmd, _, err := rootCmd.Find([]string{"server"})
	if err != nil {
		slog.Warn("config fallback", "err", err)
		return
	}

	helper.CheckFatal(viper.BindPFlag("backend.admin.server.host", serverCmd.Flags().Lookup("admin-server-host")))
	helper.CheckFatal(viper.BindPFlag("backend.admin.server.port", serverCmd.Flags().Lookup("admin-server-port")))
	helper.CheckFatal(viper.BindPFlag("backend.admin.auth.access_key", serverCmd.Flags().Lookup("admin-auth-access-key")))
	helper.CheckFatal(viper.BindPFlag("backend.admin.auth.secret_key", serverCmd.Flags().Lookup("admin-auth-secret-key")))

	helper.CheckFatal(viper.BindPFlag("backend.s3.server.host", serverCmd.Flags().Lookup("s3-server-host")))
	helper.CheckFatal(viper.BindPFlag("backend.s3.server.port", serverCmd.Flags().Lookup("s3-server-port")))
	helper.CheckFatal(viper.BindPFlag("backend.s3.storage.path", serverCmd.Flags().Lookup("s3-storage-path")))
	helper.CheckFatal(viper.BindPFlag("backend.s3.perf.buffer_size", serverCmd.Flags().Lookup("s3-perf-buffer-size")))
}

func readConfig() {
	slog.Debug("reading config")

	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("config not found, using defaults")
	} else {
		configPath := viper.GetViper().ConfigFileUsed()
		slog.Info("config read", "path", helper.TildePath(configPath))
	}
}

func parseConfig() {
	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		slog.Error("config parsing", "err", err)
		return
	}

	Admin = config.Backend.Admin
	Admin.Debug()

	IAM = config.Backend.IAM
	IAM.Debug()

	S3 = config.Backend.S3
	S3.Debug()

	relStoragePath := helper.MustResolveToRelativePath(S3.Storage.Path)
	slog.Debug("storage path resolved", "from", S3.Storage.Path, "to", relStoragePath)
	S3.Storage.Path = relStoragePath
}
