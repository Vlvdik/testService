package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	DBHost       string `mapstructure:"db_path"`
	DBPort       string `mapstructure:"db_port"`
	DBName       string `mapstructure:"db_name"`
	DBUser       string `mapstructure:"db_user"`
	DBPwd        string `mapstructure:"db_pwd"`
	StreamName   string `mapstructure:"stream_name"`
	NatsUsername string `mapstructure:"nats_username"`
	NatsPwd      string `mapstructure:"nats_pwd"`
	NatsAddr     string `mapstructure:"nats_addr"`
}

func NewConfig() (*Config, error) {
	viper.AddConfigPath("../../internal/config")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("Config", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
