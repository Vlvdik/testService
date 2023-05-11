package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server
	DB
	Cache
	Broker
	ClickHouse
}

type Server struct {
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	RequestTimeout int    `mapstructure:"request_timeout"`
}

type DB struct {
	Port string `mapstructure:"port"`
	Name string `mapstructure:"name"`
	User string `mapstructure:"user"`
	Pwd  string `mapstructure:"pwd"`
}

type Cache struct {
	Host          string `mapstructure:"host"`
	Port          string `mapstructure:"port"`
	Pwd           string `mapstructure:"pwd"`
	Secret        string `mapstructure:"secret"`
	WriteDuration int    `mapstructure:"write_duration"`
}

type Broker struct {
	URL        string `mapstructure:"url"`
	User       string `mapstructure:"user"`
	Pwd        string `mapstructure:"pwd"`
	Addr       string `mapstructure:"addr"`
	Subject    string `mapstructure:"subject"`
	MaxPending int    `mapstructure:"max_pending"`
}

type ClickHouse struct {
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
	DB      string `mapstructure:"db"`
	User    string `mapstructure:"user"`
	Pwd     string `mapstructure:"pwd"`
	Table   string `mapstructure:"table"`
	Timeout int    `mapstructure:"timeout"`
}

func NewConfig() (*Config, error) {
	viper.AddConfigPath("./src/internal/config")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("server", &cfg.Server)
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("db", &cfg.DB)
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("cache", &cfg.Cache)
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("broker", &cfg.Broker)
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("clickhouse", &cfg.ClickHouse)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
