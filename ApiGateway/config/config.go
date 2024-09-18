package config

import "github.com/spf13/viper"

type Config struct {
	Server struct {
		Port string
	}
	Services struct {
		AuthService  string
		UserService  string
		IeltsService string
	}
}

func LoadConfig() (*Config, error) {

	viper.SetConfigFile("config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
