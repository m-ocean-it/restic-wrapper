package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"restic-wrapper/secrets"

	"gopkg.in/yaml.v3"
)

var CONFIG_PATH_ENV_VAR = "RESTIC_WRAPPER_CONFIG_PATH"
var DEFAULT_CONFIG_PATH = "/etc/restic-wrapper/config.yaml"
var DEFAULT_SECRETS_PATH = "/etc/restic-wrapper/secrets.yaml"

type config struct {
	Profiles map[string]profile
}

func (c config) validate() {
	for _, prof := range c.Profiles {
		prof.validate()
	}
}

type profile struct {
	SecretsFilePath string `yaml:"secrets-file"`

	Endpoint  string `yaml:"endpoint"`
	Bucket    string `yaml:"bucket"`
	BucketDir string `yaml:"bucket-dir"`

	Sources []string `yaml:"sources"`
}

func Build() (config, error) {
	config_file_path := os.Getenv(CONFIG_PATH_ENV_VAR)

	if len(config_file_path) == 0 {
		log.Printf("Path to the config file was not specified in the %s environment variable, therefore using the default path: %s\n",
			CONFIG_PATH_ENV_VAR,
			DEFAULT_CONFIG_PATH)
		config_file_path = DEFAULT_CONFIG_PATH
	}

	var conf config

	data, err := os.ReadFile(config_file_path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = fmt.Errorf("config file does not exist at: %s", config_file_path)
			return conf, err
		}
		return conf, err
	}

	err = yaml.Unmarshal(data, &conf)

	conf.validate()

	return conf, err
}

func (c profile) Url() string {
	return fmt.Sprintf("https://%s/%s/%s",
		c.Endpoint,
		c.Bucket,
		c.BucketDir)
}

func (c profile) Secrets() secrets.Secrets {
	return secrets.Build(c.SecretsFilePath)
}

func (c *profile) validate() {
	if len(c.BucketDir) == 0 {
		log.Println("config: BUCKET_DIR empty, therefore using the bucket's root")
	}

	if len(c.SecretsFilePath) == 0 {
		log.Printf("config: SECRETS_FILE empty, therefore using the default: %s\n",
			DEFAULT_SECRETS_PATH)
		c.SecretsFilePath = DEFAULT_SECRETS_PATH
	}

	{
		err_msg := "config error"

		if len(c.Endpoint) == 0 {
			log.Fatalf("%s: ENDPOINT empty\n", err_msg)
		}
		if len(c.Bucket) == 0 {
			log.Fatalf("%s: BUCKET empty\n", err_msg)
		}
		if len(c.Sources) == 0 {
			log.Fatalf("%s: BACKUP_PATHS empty\n", err_msg)
		}
	}
}
