package credentials

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/cli/internal/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"go.yaml.in/yaml/v3"
)

const (
	CredentialsFilename = "credentials.yml"

	AccessKeyEnvVar = "LABSTORE_ACCESS_KEY"
	SecretKeyEnvVar = "LABSTORE_SECRET_KEY"
)

type Credentials struct {
	DefaultProfile *string             `yaml:"default_profile"`
	Profiles       map[string]*Profile `yaml:"profiles"`
}

type Profile struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}

func Init() {
	credentialsPath := helper.Must(CredentialsPath())

	if !helper.FileExists(credentialsPath) {
		slog.Info("creating empty credentials file")
		credentials := Credentials{
			DefaultProfile: helper.StringPtr("default"),
			Profiles:       map[string]*Profile{},
		}
		helper.CheckFatal(credentials.Save())
	}
}

func LoadCredentials() (*Credentials, error) {
	slog.Debug("loading credentials")

	credentialsPath := helper.Must(CredentialsPath())

	file, err := os.Open(credentialsPath)
	if err != nil {
		return nil, err
	}

	dec := yaml.NewDecoder(file)
	dec.KnownFields(true)

	var credentials Credentials
	if err := dec.Decode(&credentials); err != nil {
		return nil, err
	}

	return &credentials, nil
}

func (c *Credentials) Save() error {
	slog.Info("saving credentials")

	credentialsPath := helper.Must(CredentialsPath())
	credentialsDir := filepath.Dir(credentialsPath)

	if err := os.MkdirAll(credentialsDir, 0o700); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(credentialsPath, data, 0o600)
	if err != nil {
		return err
	}

	return nil
}

func LoadProfile(name string) (*Profile, error) {
	slog.Info("loading profile", "name", name)

	credentials, err := LoadCredentials()
	if err != nil {
		return nil, err
	}

	profile, ok := credentials.Profiles[name]
	if !ok {
		return nil, &errs.ErrProfileNotFound{Name: name}
	}

	return profile, nil
}

func LoadDefaultProfile() (*Profile, error) {
	slog.Info("loading default profile")

	accessKey, accessKeyOK := os.LookupEnv(AccessKeyEnvVar)
	secretKey, secretKeyOK := os.LookupEnv(SecretKeyEnvVar)

	if accessKeyOK && !secretKeyOK || !accessKeyOK && secretKeyOK {
		slog.Warn(
			"load default profile: ignoring incomplete environment",
			AccessKeyEnvVar, accessKey,
			SecretKeyEnvVar, secretKey,
		)
	}

	if accessKeyOK && secretKeyOK {
		slog.Debug("load default profile: default profile with environment credentials")
		profile := &Profile{
			AccessKey: accessKey,
			SecretKey: secretKey,
		}
		return profile, nil
	}

	credentials, err := LoadCredentials()
	if err != nil {
		return nil, err
	}

	if credentials.DefaultProfile == nil {
		return nil, &errs.ErrDefaultProfileNotSet{}
	}

	profile, ok := credentials.Profiles[*credentials.DefaultProfile]
	if !ok {
		return nil, &errs.ErrProfileNotFound{Name: *credentials.DefaultProfile}
	}

	slog.Debug("default profile from credentials file")
	return profile, nil
}

func CredentialsPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(configDir, config.Basename, CredentialsFilename)
	return configPath, nil
}
