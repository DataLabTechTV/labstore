package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/IllumiKnowLabs/labstore/server/pkg/constants"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/security"
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
	DefaultAdminSecretKey        string = DefaultAdminAuthAccessKey
	DisplayDefaultAdminSecretKey bool   = true
)

type AppConfig struct {
	Server *ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
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
	Address *AddressConfig `mapstructure:"address"`
	Auth    *AuthConfig    `mapstructure:"auth"`
}

type IAMConfig struct {
	Address *AddressConfig `mapstructure:"address"`
	DB      *IAMDBConfig   `mapstructure:"db"`
}

type S3Config struct {
	Address *AddressConfig `mapstructure:"address"`
	Paging  *PagingConfig  `mapstructure:"paging"`
	IO      *IOConfig      `mapstructure:"io"`
}

type AddressConfig struct {
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

var App *AppConfig

func (config *AppConfig) Debug() {
	config.Server.Debug()
}

func (config *ServerConfig) Debug() {
	config.Storage.Debug()
	config.Admin.Debug()
	config.IAM.Debug()
	config.S3.Debug()
}

func (config *StorageConfig) Debug() {
	slog.Debug("config set", "name", "server.storage.data_dir", "value", config.DataDir)
	slog.Debug("config set", "name", "server.storage.keys_dir", "value", config.KeysDir)
}

func (config *StorageConfig) PreparePaths() {
	relStorageDataDir := helper.MustRelPath(config.DataDir)
	slog.Debug("storage data dir resolved", "from", config.DataDir, "to", relStorageDataDir)
	config.DataDir = relStorageDataDir

	config.ObjectsPath = filepath.Join(config.DataDir, DefaultStorageObjectsDir)
	slog.Debug("object storage path set", "path", config.ObjectsPath)

	config.MetadataPath = filepath.Join(config.DataDir, DefaultStorageMetadataDir)
	slog.Debug("metadata storage path set", "path", config.MetadataPath)

	relStorageKeysDir := helper.MustRelPath(config.KeysDir)
	slog.Debug("storage keys dir resolved", "from", config.KeysDir, "to", relStorageKeysDir)
	config.KeysDir = relStorageKeysDir

	config.MasterKeyPath = filepath.Join(config.KeysDir, DefaultMasterKeyFilename)
	slog.Debug("master key path set", "path", config.MasterKeyPath)
}

func (config *AdminConfig) Debug() {
	slog.Debug("config set", "name", "server.admin.address.host", "value", config.Address.Host)
	slog.Debug("config set", "name", "server.admin.address.port", "value", config.Address.Port)
	slog.Debug("config set", "name", "server.admin.access_key", "value", config.Auth.AccessKey)

	var adminSecretKeyDisplay string
	if len(config.Auth.SecretKey) > 0 {
		adminSecretKeyDisplay = security.Redacted
	} else {
		adminSecretKeyDisplay = constants.Empty
	}
	slog.Debug("config set", "name", "server.admin.secret_key", "value", adminSecretKeyDisplay)
}

func (config *IAMConfig) Debug() {
	slog.Debug("config set", "name", "server.iam.address.host", "value", config.Address.Host)
	slog.Debug("config set", "name", "server.iam.address.port", "value", config.Address.Port)
	slog.Debug("config set", "name", "server.iam.db.max_open_conns", "value", config.DB.MaxOpenConns)
	slog.Debug("config set", "name", "server.iam.db.max_idle_conns", "value", config.DB.MaxIdleConns)
	slog.Debug("config set", "name", "server.iam.db.write_chan_cap", "value", config.DB.WriteChanCap)
	slog.Debug("config set", "name", "server.iam.db.timeout_ms", "value", config.DB.TimeoutMs)
	slog.Debug("config set", "name", "server.iam.db.read_cache_size_kib", "value", config.DB.ReadCacheSizeKiB)
	slog.Debug("config set", "name", "server.iam.db.write_cache_size_kib", "value", config.DB.WriteCacheSizeKiB)
}

func (config *S3Config) Debug() {
	slog.Debug("config set", "name", "server.s3.address.host", "value", config.Address.Host)
	slog.Debug("config set", "name", "server.s3.address.port", "value", config.Address.Port)
	slog.Debug("config set", "name", "server.s3.paging.max_keys", "value", config.Paging.MaxKeys)
	slog.Debug("config set", "name", "server.s3.io.buffer_size", "value", config.IO.BufferSize)
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
	viper.SetDefault("server.storage.data_dir", DefaultStorageDataDir)
	viper.SetDefault("server.storage.keys_dir", DefaultStorageKeysDir)

	// admin
	viper.SetDefault("server.admin.address.host", DefaultAdminServerHost)
	viper.SetDefault("server.admin.address.port", DefaultAdminServerPort)
	viper.SetDefault("server.admin.auth.access_key", DefaultAdminAuthAccessKey)

	randomSecretKey, err := security.GeneratePassword(32)
	if err != nil {
		slog.Warn("admin password generation", "err", err)
	} else {
		DefaultAdminSecretKey = randomSecretKey
	}

	viper.SetDefault("server.admin.auth.secret_key", DefaultAdminSecretKey)

	// iam
	viper.SetDefault("server.iam.address.host", DefaultIAMServerHost)
	viper.SetDefault("server.iam.address.port", DefaultIAMServerPort)
	viper.SetDefault("server.iam.db.max_open_conns", DefaultIAMDBMaxOpenConns)
	viper.SetDefault("server.iam.db.max_idle_conns", DefaultIAMDBMaxIdleConns)
	viper.SetDefault("server.iam.db.write_chan_cap", DefaultIAMWriteChanCap)
	viper.SetDefault("server.iam.db.timeout_ms", DefaultIAMDBTimeoutMs)
	viper.SetDefault("server.iam.db.read_cache_size_kib", DefaultIAMDBReadCacheSizeKiB)
	viper.SetDefault("server.iam.db.write_cache_size_kib", DefaultIAMDBWriteCacheSizeKiB)

	// s3
	viper.SetDefault("server.s3.address.host", DefaultS3ServerHost)
	viper.SetDefault("server.s3.address.port", DefaultS3ServerPort)
	viper.SetDefault("server.s3.paging.max_keys", DefaultS3PagingMaxKeys)
	viper.SetDefault("server.s3.io.buffer_size", DefaultS3IOBufferSize)
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
	bindPFlagIfExists("server.storage.data_dir", serverCmd, "storage-data-dir")
	bindPFlagIfExists("server.storage.keys_dir", serverCmd, "storage-keys-dir")

	// admin
	bindPFlagIfExists("server.admin.address.host", serverCmd, "admin-server-host")
	bindPFlagIfExists("server.admin.address.port", serverCmd, "admin-server-port")
	bindPFlagIfExists("server.admin.auth.access_key", serverCmd, "admin-auth-access-key")
	bindPFlagIfExists("server.admin.auth.secret_key", serverCmd, "admin-auth-secret-key")

	// iam
	bindPFlagIfExists("server.iam.address.host", serverCmd, "iam-server-host")
	bindPFlagIfExists("server.aim.address.port", serverCmd, "iam-server-port")
	bindPFlagIfExists("server.iam.db.max_open_conns", serverCmd, "iam-db-max-open-conns")
	bindPFlagIfExists("server.iam.db.max_idle_conns", serverCmd, "iam-db-max-idle-conns")
	bindPFlagIfExists("server.iam.db.write_chan_cap", serverCmd, "iam-db-write-chan-cap")
	bindPFlagIfExists("server.iam.db.timeout_ms", serverCmd, "iam-db-timeout-ms")
	bindPFlagIfExists("server.iam.db.read_cache_size_kib", serverCmd, "iam-db-read-cache-size-kib")
	bindPFlagIfExists("server.iam.db.write_cache_size_kib", serverCmd, "iam-db-write-cache-size-kib")

	// s3
	bindPFlagIfExists("server.s3.address.host", serverCmd, "s3-server-host")
	bindPFlagIfExists("server.s3.address.port", serverCmd, "s3-server-port")
	bindPFlagIfExists("server.s3.paging.max_keys", serverCmd, "s3-paging-max-keys")
	bindPFlagIfExists("server.s3.io.buffer_size", serverCmd, "s3-io-buffer-size")
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
	App = new(AppConfig)
	if err := viper.Unmarshal(&App); err != nil {
		slog.Error("config parsing", "err", err)
		return
	}

	App.Debug()
	App.Server.Storage.PreparePaths()

	if DisplayDefaultAdminSecretKey && App.Server.Admin.Auth.SecretKey == DefaultAdminSecretKey {
		slog.Warn("no secret ket set for admin, randomly generating")
		fmt.Printf(
			"🔑 Temporary secret key for %s: %s\n",
			App.Server.Admin.Auth.AccessKey, DefaultAdminSecretKey,
		)
	}
}
