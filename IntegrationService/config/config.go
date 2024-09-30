package config

import (
	"os"

	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
