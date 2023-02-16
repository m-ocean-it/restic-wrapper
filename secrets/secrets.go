package secrets

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Secrets struct {
	Aws            AwsCredentials `yaml:"aws"`
	ResticPassword string         `yaml:"restic_password"`
}

type AwsCredentials struct {
	KeyId string `yaml:"key_id"`
	Key   string `yaml:"key"`
}

func Build(filepath string) Secrets {
	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalln(err)
	}

	var cred Secrets
	err = yaml.Unmarshal(file, &cred)
	if err != nil {
		log.Fatalln(err)
	}

	cred.validate()

	return cred
}

func (c Secrets) validate() {
	if len(c.ResticPassword) == 0 {
		panic("restic password is empty")
	}
	if len(c.Aws.KeyId) == 0 {
		panic("aws key id is empty")
	}
	if len(c.Aws.Key) == 0 {
		panic("aws key is empty")
	}
}
