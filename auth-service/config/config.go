package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Redis struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"redis"`
	Telegram struct {
		BotToken string `yaml:"bot_token"`
	} `yaml:"telegram"`
	GRPC struct {
		UserService struct {
			Address string `yaml:"address"`
		} `yaml:"userService"`
	} `yaml:"grpc"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DbName   string `yaml:"dbname"`
		Password string `yaml:"password"`
		Username string `yaml:"username"`
	} `yaml:"database"`
}

func LoadConfig() (*Config, error) {
	f, err := os.Open("config/config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
