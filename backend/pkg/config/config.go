package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/constants"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/security"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Basename = "labstore"

	DefaultStorageDataDir     = "./data"
	DefaultStorageObjectsDir  = "objects"
	DefaultStorageMetadataDir = "metadata"
	DefaultStorageKeysDir     = "./keys"
	DefaultMasterKeyFilename  = "master.key"

	DefaultAdminServerHost    = "0.0.0.0"
	DefaultAdminServerPort    = 6787
	DefaultAdminAuthAccessKey = "admin"

	DefaultIAMServerHost     = "0.0.0.0"
	DefaultIAMServerPort     = 6788
	DefaultIAMDBMaxOpenConns = 3
	DefaultIAMDBMaxIdleConns = 3
	DefaultIAMWriteChanCap   = 32

	DefaultIAMDBTimeoutMs         = 5000
	DefaultIAMDBReadCacheSizeKiB  = 65536
	DefaultIAMDBWriteCacheSizeKiB = 16384

	DefaultS3ServerHost    = "0.0.0.0"
	DefaultS3ServerPort    = 6789
	DefaultS3PagingMaxKeys = 1000
	DefaultS3IOBufferSize  = 256 * helper.KiB
)

var (
	DefaultAdminSecretKey string = DefaultAdminAuthAccessKey
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
	Paging *PagingConfig `mapstructure:"paging"`
	IO     *IOConfig     `mapstructure:"io"`
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
	WriteChanCap      int `mapstructure:"write_chan_cap"`
	TimeoutMs         int `mapstructure:"timeout_ms"`
	ReadCacheSizeKiB  int `mapstructure:"read_cache_size_kib"`
	WriteCacheSizeKiB int `mapstructure:"write_cache_size_kib"`
}

type PagingConfig struct {
	MaxKeys int `mapstructure:"max_keys"`
}

type IOConfig struct {
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
	slog.Debug("config set", "name", "backend.iam.db.write_chan_cap", "value", config.DB.WriteChanCap)
	slog.Debug("config set", "name", "backend.iam.db.timeout_ms", "value", config.DB.TimeoutMs)
	slog.Debug("config set", "name", "backend.iam.db.read_cache_size_kib", "value", config.DB.ReadCacheSizeKiB)
	slog.Debug("config set", "name", "backend.iam.db.write_cache_size_kib", "value", config.DB.WriteCacheSizeKiB)
}

func (config *S3Config) Debug() {
	slog.Debug("config set", "name", "backend.s3.server.host", "value", config.Server.Host)
	slog.Debug("config set", "name", "backend.s3.server.port", "value", config.Server.Port)
	slog.Debug("config set", "name", "backend.s3.paging.max_keys", "value", config.Paging.MaxKeys)
	slog.Debug("config set", "name", "backend.s3.io.buffer_size", "value", config.IO.BufferSize)
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

	viper.SetConfigName(Basename)

	viper.AddConfigPath(".")

	systemConfigPath := filepath.Join("/etc", Basename)
	viper.AddConfigPath(systemConfigPath)

	if configDir, err := os.UserConfigDir(); err != nil {
		slog.Warn("user config dir skipped", "err", err)
	} else {
		userConfigPath := filepath.Join(configDir, Basename)
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

	randomSecretKey, err := security.GeneratePassword(32)
	if err != nil {
		slog.Warn("admin password generation", "err", err)
	} else {
		DefaultAdminSecretKey = randomSecretKey
	}

	viper.SetDefault("backend.admin.auth.secret_key", DefaultAdminSecretKey)

	// iam
	viper.SetDefault("backend.iam.server.host", DefaultIAMServerHost)
	viper.SetDefault("backend.iam.server.port", DefaultIAMServerPort)
	viper.SetDefault("backend.iam.db.max_open_conns", DefaultIAMDBMaxOpenConns)
	viper.SetDefault("backend.iam.db.max_idle_conns", DefaultIAMDBMaxIdleConns)
	viper.SetDefault("backend.iam.db.write_chan_cap", DefaultIAMWriteChanCap)
	viper.SetDefault("backend.iam.db.timeout_ms", DefaultIAMDBTimeoutMs)
	viper.SetDefault("backend.iam.db.read_cache_size_kib", DefaultIAMDBReadCacheSizeKiB)
	viper.SetDefault("backend.iam.db.write_cache_size_kib", DefaultIAMDBWriteCacheSizeKiB)

	// s3
	viper.SetDefault("backend.s3.server.host", DefaultS3ServerHost)
	viper.SetDefault("backend.s3.server.port", DefaultS3ServerPort)
	viper.SetDefault("backend.s3.paging.max_keys", DefaultS3PagingMaxKeys)
	viper.SetDefault("backend.s3.io.buffer_size", DefaultS3IOBufferSize)
}

func setOverrides(rootCmd *cobra.Command) {
	slog.Debug("setting config env and cli overrides")

	viper.SetEnvPrefix(Basename)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	if rootCmd == nil {
		slog.Warn("setting overrides without root command")
		return
	}

	serverCmd, _, err := rootCmd.Find([]string{"server"})
	if err != nil {
		slog.Warn("config fallback", "err", err)
		return
	}

	// storage
	bindPFlagIfExists("backend.storage.data_dir", serverCmd, "storage-data-dir")
	bindPFlagIfExists("backend.storage.keys_dir", serverCmd, "storage-keys-dir")

	// admin
	bindPFlagIfExists("backend.admin.server.host", serverCmd, "admin-server-host")
	bindPFlagIfExists("backend.admin.server.port", serverCmd, "admin-server-port")
	bindPFlagIfExists("backend.admin.auth.access_key", serverCmd, "admin-auth-access-key")
	bindPFlagIfExists("backend.admin.auth.secret_key", serverCmd, "admin-auth-secret-key")

	// iam
	bindPFlagIfExists("backend.iam.server.host", serverCmd, "iam-server-host")
	bindPFlagIfExists("backend.aim.server.port", serverCmd, "iam-server-port")
	bindPFlagIfExists("backend.iam.db.max_open_conns", serverCmd, "iam-db-max-open-conns")
	bindPFlagIfExists("backend.iam.db.max_idle_conns", serverCmd, "iam-db-max-idle-conns")
	bindPFlagIfExists("backend.iam.db.write_chan_cap", serverCmd, "iam-db-write-chan-cap")
	bindPFlagIfExists("backend.iam.db.timeout_ms", serverCmd, "iam-db-timeout-ms")
	bindPFlagIfExists("backend.iam.db.read_cache_size_kib", serverCmd, "iam-db-read-cache-size-kib")
	bindPFlagIfExists("backend.iam.db.write_cache_size_kib", serverCmd, "iam-db-write-cache-size-kib")

	// s3
	bindPFlagIfExists("backend.s3.server.host", serverCmd, "s3-server-host")
	bindPFlagIfExists("backend.s3.server.port", serverCmd, "s3-server-port")
	bindPFlagIfExists("backend.s3.paging.max_keys", serverCmd, "s3-paging-max-keys")
	bindPFlagIfExists("backend.s3.io.buffer_size", serverCmd, "s3-io-buffer-size")
}

func bindPFlagIfExists(configKey string, cmd *cobra.Command, flagName string) {
	flag := cmd.Flags().Lookup(flagName)
	if flag != nil {
		helper.CheckFatal(viper.BindPFlag(configKey, flag))
	}
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

	if Admin.Auth.SecretKey == DefaultAdminSecretKey {
		slog.Warn("no secret ket set for admin, randomly generating")
		fmt.Printf(
			"🔑 Temporary secret key for %s: %s\n",
			Admin.Auth.AccessKey, DefaultAdminSecretKey,
		)
	}
}
