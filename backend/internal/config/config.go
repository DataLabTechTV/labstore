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

const configBasename = "labstore"

const DefaultHost = "0.0.0.0"
const DefaultPort = 6789
const DefaultStoragePath = "./data"
const DefaultAdminAccessKey = "admin"
const DefaultAdminSecretKey = DefaultAdminAccessKey
const DefaultPerfBufferSize = 256 * helper.KiB

var Config AppConfig

type AppConfig struct {
	Server ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Host        string            `mapstructure:"host"`
	Port        uint16            `mapstructure:"port"`
	Storage     StorageConfig     `mapstructure:"storage"`
	Admin       AdminConfig       `mapstructure:"admin"`
	Performance PerformanceConfig `mapstructure:"performance"`
}

type StorageConfig struct {
	Path string `mapstructure:"path"`
}

type AdminConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type PerformanceConfig struct {
	BufferSize int `mapstructure:"buffer_size"`
}

func (config *ServerConfig) Debug() {
	slog.Debug("config set", "name", "server.host", "value", config.Host)
	slog.Debug("config set", "name", "server.port", "value", config.Port)
	slog.Debug("config set", "name", "server.storage.path", "value", config.Storage.Path)
	slog.Debug("config set", "name", "server.admin.access_key", "value", config.Admin.AccessKey)
	slog.Debug("config set", "name", "server.performance.buffer_size", "value", config.Performance.BufferSize)

	var adminSecretKeyDisplay string
	if len(config.Admin.SecretKey) > 0 {
		adminSecretKeyDisplay = security.Redacted
	} else {
		adminSecretKeyDisplay = constants.Empty
	}
	slog.Debug("config set", "name", "server.admin.secret_key", "value", adminSecretKeyDisplay)
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

	viper.SetDefault("server.host", DefaultHost)
	viper.SetDefault("server.port", DefaultPort)
	viper.SetDefault("server.storage.path", DefaultStoragePath)
	viper.SetDefault("server.admin.access_key", DefaultAdminAccessKey)

	defaultAdminSecretKey, err := security.GeneratePassword(32)
	if err != nil {
		slog.Error("admin password generation", "err", err)
		defaultAdminSecretKey = DefaultAdminSecretKey
	}

	viper.SetDefault("server.admin.secret_key", defaultAdminSecretKey)
	fmt.Printf("🔑 Default admin secret key: %s\n", defaultAdminSecretKey)

	viper.SetDefault("server.performance.buffer_size", DefaultPerfBufferSize)
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

	helper.CheckFatal(viper.BindPFlag("server.host", serverCmd.Flags().Lookup("host")))
	helper.CheckFatal(viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port")))
	helper.CheckFatal(viper.BindPFlag("server.storage.path", serverCmd.Flags().Lookup("storage-path")))
	helper.CheckFatal(viper.BindPFlag("server.admin.access_key", serverCmd.Flags().Lookup("admin-user")))
	helper.CheckFatal(viper.BindPFlag("server.admin.secret_key", serverCmd.Flags().Lookup("admin-pass")))
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
	if err := viper.Unmarshal(&Config); err != nil {
		slog.Error("config parsing", "err", err)
		return
	}

	Config.Server.Debug()

	relStoragePath := helper.MustResolveToRelativePath(Config.Server.Storage.Path)
	slog.Debug("storage path resolved", "from", Config.Server.Storage.Path, "to", relStoragePath)
	Config.Server.Storage.Path = relStoragePath
}
