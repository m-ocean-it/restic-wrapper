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
	SecretsFilePath string `yaml:"SECRETS_FILE"`

	ENDPOINT   string `yaml:"ENDPOINT"`
	BUCKET     string `yaml:"BUCKET"`
	BUCKET_DIR string `yaml:"BUCKET_DIR"`

	BACKUP_PATHS []string `yaml:"BACKUP_PATHS"`
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
	data, err := os.ReadFile(config_file_path) // TODO: also, use env var
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

func (c config) Url() string {
	return fmt.Sprintf("https://%s/%s/%s",
		c.ENDPOINT,
		c.BUCKET,
		c.BUCKET_DIR)
}

func (c config) Secrets() secrets.Secrets {
	return secrets.Build(c.SecretsFilePath)
}

func (c *config) validate() {
	if len(c.BUCKET_DIR) == 0 {
		log.Println("config: BUCKET_DIR empty, therefore using the bucket's root")
	}

	if len(c.SecretsFilePath) == 0 {
		log.Printf("config: SECRETS_FILE empty, therefore using the default: %s\n",
			DEFAULT_SECRETS_PATH)
		c.SecretsFilePath = DEFAULT_SECRETS_PATH
	}

	{
		err_msg := "config error"

		if len(c.ENDPOINT) == 0 {
			log.Fatalf("%s: ENDPOINT empty\n", err_msg)
		}
		if len(c.BUCKET) == 0 {
			log.Fatalf("%s: BUCKET empty\n", err_msg)
		}
		if len(c.BACKUP_PATHS) == 0 {
			log.Fatalf("%s: BACKUP_PATHS empty\n", err_msg)
		}
	}
}
