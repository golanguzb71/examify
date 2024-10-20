package config

import (
	"os"

	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Grpc     GrpcConfig     `yaml:"grpc"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type GrpcConfig struct {
	UserService struct {
		Address string `yaml:"address"`
	} `yaml:"userService"`
	IntegrationService struct {
		Address string `yaml:"address"`
	} `yaml:"integrationService"`
	BonusService struct {
		Address string `yaml:"address"`
	} `yaml:"bonusService"`
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
