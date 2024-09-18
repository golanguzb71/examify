package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ServerConfig struct {
	Port string `yaml:"port"`
}

type GrpcConfig struct {
	IeltsService struct {
		Address string `yaml:"address"`
	} `yaml:"ieltsService"`
	UserService struct {
		Address string `yaml:"address"`
	} `yaml:"userService"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Grpc   GrpcConfig   `yaml:"grpc"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
