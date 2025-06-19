package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Emulate bool   `yaml:"emulate"`
	CTA     CTA    `yaml:"cta"`
	NYCMTA  NYCMTA `yaml:"nycmta"`
}

type NYCMTA struct {
	APIKey string  `yaml:"api_key"`
	Bus    MTAInfo `yaml:"bus"`
	Train  MTAInfo `yaml:"train"`
}

type CTA struct {
	Bus   CTAInfo `yaml:"bus"`
	Train CTAInfo `yaml:"train"`
}

type MTAInfo struct {
	StopIDs []string `yaml:"stop_ids"`
}

type CTAInfo struct {
	StopIDs []int  `yaml:"stop_ids"`
	APIKey  string `yaml:"api_key"`
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
