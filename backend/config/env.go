package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/DataLabTechTV/labstore/backend/internal/helper"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const dotenvPath = ".env"

var Env ServerConfig

type ServerConfig struct {
	Port           uint16 `env:"LS_PORT" envDefault:"6789"`
	StorageRoot    string `env:"LS_STORAGE_ROOT" envDefault:"../data"`
	AdminAccessKey string `env:"LS_ADMIN_ACCESS_KEY" envDefault:"admin"`
	AdminSecretKey string `env:"LS_ADMIN_SECRET_KEY" envDefault:"admin"`
}

func LoadEnv() {
	if err := godotenv.Load(dotenvPath); err != nil {
		log.Debug("No .env file found, skipping...")
	} else {
		log.Debugln("Environment source:", dotenvPath)
	}

	Env = helper.Must(env.ParseAs[ServerConfig]())

	cwd := helper.Must(os.Getwd())
	absStoragePath := helper.Must(filepath.Abs(Env.StorageRoot))
	Env.StorageRoot = helper.Must(filepath.Rel(cwd, absStoragePath))

	t := reflect.TypeOf(Env)
	v := reflect.ValueOf(Env)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		env_var_name := field.Tag.Get("env")
		env_var_value := fmt.Sprintf("%v", value)

		if strings.Contains(env_var_name, "SECRET") {
			if len(env_var_value) > 0 {
				env_var_value = "[REDACTED]"
			} else {
				env_var_value = "[EMPTY]"
			}
		}

		log.Debugf("%s: %s", env_var_name, env_var_value)
	}
}
