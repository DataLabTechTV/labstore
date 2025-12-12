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

	DefaultStorageDataDir     = "./data"
	DefaultStorageObjectsDir  = "objects"
	DefaultStorageMetadataDir = "metadata"
	DefaultStorageKeysDir     = "./keys"
	DefaultMasterKeyFilename  = "master.key"

	DefaultAdminServerHost    = "0.0.0.0"
	DefaultAdminServerPort    = 6787
	DefaultAdminAuthAccessKey = "admin"
	DefaultAdminSecretKey     = DefaultAdminAuthAccessKey

	DefaultIAMServerHost          = "0.0.0.0"
	DefaultIAMServerPort          = 6788
	DefaultIAMDBMaxOpenConns      = 3
	DefaultIAMDBMaxIdleConns      = 3
	DefaultIAMDBTimeoutMs         = 5000
	DefaultIAMDBReadCacheSizeKiB  = 65536
	DefaultIAMDBWriteCacheSizeKiB = 16384

	DefaultS3ServerHost     = "0.0.0.0"
	DefaultS3ServerPort     = 6789
	DefaultS3PerfBufferSize = 256 * helper.KiB
)

type AppConfig struct {
	Backend *BackendConfig `mapstructure:"backend"`
}

type BackendConfig struct {
	Storage *StorageConfig `mapstructure:"storage"`
	Admin   *AdminConfig   `mapstructure:"admin"`
	IAM     *IAMConfig     `mapstructure:"iam"`
	S3      *S3Config      `mapstructure:"s3"`
}

type StorageConfig struct {
	DataDir      string `mapstructure:"data_dir"`
	ObjectsPath  string `mapstructure:"-"`
	MetadataPath string `mapstructure:"-"`

	KeysDir       string `mapstructure:"keys_dir"`
	MasterKeyPath string `mapstructure:"-"`
}

type AdminConfig struct {
	Server *ServerConfig `mapstructure:"server"`
	Auth   *AuthConfig   `mapstructure:"auth"`
}

type IAMConfig struct {
	Server *ServerConfig `mapstructure:"server"`
	DB     *IAMDBConfig  `mapstructure:"db"`
}

type S3Config struct {
	Server *ServerConfig `mapstructure:"server"`
	Perf   *PerfConfig   `mapstructure:"perf"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

type AuthConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type IAMDBConfig struct {
	MaxOpenConns      int `mapstructure:"max_open_conns"`
	MaxIdleConns      int `mapstructure:"max_idle_conns"`
	TimeoutMs         int `mapstructure:"timeout_ms"`
	ReadCacheSizeKiB  int `mapstructure:"read_cache_size_kib"`
	WriteCacheSizeKiB int `mapstructure:"write_cache_size_kib"`
}

type PerfConfig struct {
	BufferSize int `mapstructure:"buffer_size"`
}

var Storage *StorageConfig
var Admin *AdminConfig
var IAM *IAMConfig
var S3 *S3Config

func (config *StorageConfig) Debug() {
	slog.Debug("config set", "name", "backend.storage.data_dir", "value", config.DataDir)
	slog.Debug("config set", "name", "backend.storage.keys_dir", "value", config.KeysDir)
}

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
	slog.Debug("config set", "name", "backend.iam.db.max_open_conns", "value", config.DB.MaxOpenConns)
	slog.Debug("config set", "name", "backend.iam.db.max_idle_conns", "value", config.DB.MaxIdleConns)
	slog.Debug("config set", "name", "backend.iam.db.timeout_ms", "value", config.DB.TimeoutMs)
	slog.Debug("config set", "name", "backend.iam.db.read_cache_size_kib", "value", config.DB.ReadCacheSizeKiB)
	slog.Debug("config set", "name", "backend.iam.db.write_cache_size_kib", "value", config.DB.WriteCacheSizeKiB)
}

func (config *S3Config) Debug() {
	slog.Debug("config set", "name", "backend.s3.server.host", "value", config.Server.Host)
	slog.Debug("config set", "name", "backend.s3.server.port", "value", config.Server.Port)
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

	// storage
	viper.SetDefault("backend.storage.data_dir", DefaultStorageDataDir)
	viper.SetDefault("backend.storage.keys_dir", DefaultStorageKeysDir)

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
	viper.SetDefault("backend.iam.db.max_open_conns", DefaultIAMDBMaxOpenConns)
	viper.SetDefault("backend.iam.db.max_idle_conns", DefaultIAMDBMaxIdleConns)
	viper.SetDefault("backend.iam.db.timeout_ms", DefaultIAMDBTimeoutMs)
	viper.SetDefault("backend.iam.db.read_cache_size_kib", DefaultIAMDBReadCacheSizeKiB)
	viper.SetDefault("backend.iam.db.write_cache_size_kib", DefaultIAMDBWriteCacheSizeKiB)

	// s3
	viper.SetDefault("backend.s3.server.host", DefaultS3ServerHost)
	viper.SetDefault("backend.s3.server.port", DefaultS3ServerPort)
	viper.SetDefault("backend.s3.perf.buffer_size", DefaultS3PerfBufferSize)
}

func setOverrides(rootCmd *cobra.Command) {
	slog.Debug("setting config env and cli overrides")

	viper.SetEnvPrefix(configBasename)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	if rootCmd == nil {
		slog.Warn("set overrides without root command")
		return
	}

	serverCmd, _, err := rootCmd.Find([]string{"server"})
	if err != nil {
		slog.Warn("config fallback", "err", err)
		return
	}

	// storage
	helper.CheckFatal(viper.BindPFlag("backend.storage.data_dir", serverCmd.Flags().Lookup("storage-data-dir")))
	helper.CheckFatal(viper.BindPFlag("backend.storage.keys_dir", serverCmd.Flags().Lookup("storage-keys-dir")))

	// admin
	helper.CheckFatal(viper.BindPFlag("backend.admin.server.host", serverCmd.Flags().Lookup("admin-server-host")))
	helper.CheckFatal(viper.BindPFlag("backend.admin.server.port", serverCmd.Flags().Lookup("admin-server-port")))
	helper.CheckFatal(viper.BindPFlag("backend.admin.auth.access_key", serverCmd.Flags().Lookup("admin-auth-access-key")))
	helper.CheckFatal(viper.BindPFlag("backend.admin.auth.secret_key", serverCmd.Flags().Lookup("admin-auth-secret-key")))

	// iam
	helper.CheckFatal(viper.BindPFlag("backend.iam.server.host", serverCmd.Flags().Lookup("iam-server-host")))
	helper.CheckFatal(viper.BindPFlag("backend.aim.server.port", serverCmd.Flags().Lookup("iam-server-port")))
	helper.CheckFatal(viper.BindPFlag("backend.iam.db.max_open_conns", serverCmd.Flags().Lookup("iam-db-max-open-conns")))
	helper.CheckFatal(viper.BindPFlag("backend.iam.db.max_idle_conns", serverCmd.Flags().Lookup("iam-db-max-idle-conns")))
	helper.CheckFatal(viper.BindPFlag("backend.iam.db.timeout_ms", serverCmd.Flags().Lookup("iam-db-timeout-ms")))
	helper.CheckFatal(viper.BindPFlag("backend.iam.db.read_cache_size_kib", serverCmd.Flags().Lookup("iam-db-read-cache-size-kib")))
	helper.CheckFatal(viper.BindPFlag("backend.iam.db.write_cache_size_kib", serverCmd.Flags().Lookup("iam-db-write-cache-size-kib")))

	// s3
	helper.CheckFatal(viper.BindPFlag("backend.s3.server.host", serverCmd.Flags().Lookup("s3-server-host")))
	helper.CheckFatal(viper.BindPFlag("backend.s3.server.port", serverCmd.Flags().Lookup("s3-server-port")))
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

	Storage = config.Backend.Storage
	Storage.Debug()

	Admin = config.Backend.Admin
	Admin.Debug()

	IAM = config.Backend.IAM
	IAM.Debug()

	S3 = config.Backend.S3
	S3.Debug()

	relStorageDataDir := helper.MustResolveToRelativePath(Storage.DataDir)
	slog.Debug("storage data dir resolved", "from", Storage.DataDir, "to", relStorageDataDir)
	Storage.DataDir = relStorageDataDir

	Storage.ObjectsPath = filepath.Join(Storage.DataDir, DefaultStorageObjectsDir)
	slog.Debug("object storage path set", "path", Storage.ObjectsPath)

	Storage.MetadataPath = filepath.Join(Storage.DataDir, DefaultStorageMetadataDir)
	slog.Debug("metadata storage path set", "path", Storage.MetadataPath)

	relStorageKeysDir := helper.MustResolveToRelativePath(Storage.KeysDir)
	slog.Debug("storage keys dir resolved", "from", Storage.KeysDir, "to", relStorageKeysDir)
	Storage.KeysDir = relStorageKeysDir

	Storage.MasterKeyPath = filepath.Join(Storage.KeysDir, DefaultMasterKeyFilename)
	slog.Debug("master key path set", "path", Storage.MasterKeyPath)
}
