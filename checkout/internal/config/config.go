package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ProductService struct {
	Url           string `yaml:"url"`
	Token         string `yaml:"token"`
	RateLimit     uint32 `yaml:"rateLimit"`
	MaxConcurrent uint32 `yaml:"maxConcurrent"`
}

type ConfigStruct struct {
	Token    string `yaml:"token"`
	Services struct {
		Loms           string         `yaml:"loms"`
		ProductService ProductService `yaml:"productService"`
	} `yaml:"services"`
}

var ConfigData ConfigStruct

func Init() error {
	rawYAML, err := os.ReadFile("config.yml")
	if err != nil {
		return errors.WithMessage(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &ConfigData)
	if err != nil {
		return errors.WithMessage(err, "parsing yaml")
	}

	return nil
}
