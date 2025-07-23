package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Cert    Cert    `mapstructure:"cert"`
	Server  Server  `mapstructure:"server"`
	Redis   Redis   `mapstructure:"redis"`
	DB      DB      `mapstructure:"db"`
	Caching Caching `mapstructure:"caching"`
	Service Service `mapstructure:"service"`
}

type Server struct {
	Host     string `mapstructure:"host"`
	GrpcPort string `mapstructure:"grpc_port"`
	HttpPort string `mapstructure:"http_port"`
}

type Cert struct {
	PrivateKeyPath string `mapstructure:"private_key_path"`
	PublicKeyPath  string `mapstructure:"public_key_path"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	PWD      string `mapstructure:"pwd"`
	Database int    `mapstructure:"database"`
}

type Caching struct {
	Enable bool   `mapstructure:"enable"`
	TTL    string `mapstructure:"ttl"`
}

type DB struct {
	Mongo struct {
		DSN      string `mapstructure:"dsn"`
		Database string `mapstructure:"database"`
	} `mapstructure:"mongo"`
}

type Service struct {
	Session struct {
		AccessTokenTTL  string `mapstructure:"access_token_ttl"`
		RefreshTokenTTL string `mapstructure:"refresh_token_ttl"`
	} `mapstructure:"session"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	v := viper.NewWithOptions()
	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	v.SetConfigFile("config/config.yml")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
