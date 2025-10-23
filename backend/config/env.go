package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type ServerConfig struct {
	StorageRoot    string `env:"LS_STORAGE_ROOT" envDefault:"../data"`
	AdminAccessKey string `env:"LS_ADMIN_ACCESS_KEY" envDefault:"admin"`
	AdminSecretKey string `env:"LS_ADMIN_SECRET_KEY" envDefault:"admin"`
}

const DOTENV_PATH = ".env"

var Env ServerConfig

func LoadEnv() {
	if err := godotenv.Load(DOTENV_PATH); err != nil {
		log.Debug("No .env file found, skipping...")
	} else {
		log.Infoln("Environment loaded from:", DOTENV_PATH)
	}

	if err := env.Parse(&Env); err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	absStoragePath, err := filepath.Abs(Env.StorageRoot)
	if err != nil {
		log.Fatal(err)
	}

	relStorageRoot, err := filepath.Rel(cwd, absStoragePath)
	if err != nil {
		log.Fatal(err)
	}

	Env.StorageRoot = relStorageRoot

	t := reflect.TypeOf(Env)
	v := reflect.ValueOf(Env)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		env_var_name := field.Tag.Get("env")
		env_var_value := value.String()

		if strings.Contains(env_var_name, "SECRET") {
			if len(env_var_value) > 0 {
				env_var_value = "[REDACTED]"
			} else {
				env_var_value = "[EMPTY]"
			}
		}

		log.Infoln(env_var_name, env_var_value)
	}
}
