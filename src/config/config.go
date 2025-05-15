package config

import (
	"os"
	"gopkg.in/yaml.v2"
)

type Config struct {
	CTA CTA `yaml:"cta"`
	NYCMTA NYCMTA `yaml:"nyc_mta"`
}

type NYCMTA struct {
	APIKey string `yaml:"api_key"`
}

type CTA struct {
	Bus Bus `yaml:"bus"`
	Train Train `yaml:"train"`
}

type Bus struct {
	APIKey string `yaml:"api_key"`
}

type Train struct {
	APIKey string `yaml:"api_key"`
}

func NewConfig() (*Config, error) {
	c := Config{}

	data, err := os.ReadFile("./configs/config.yaml")
	if err != nil {
		return &c, err
	}

	err = yaml.Unmarshal([]byte(data), &c)
	return &c, err
}
