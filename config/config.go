package config

import (
	"auth-svc/internal/types"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server     Server     `mapstructure:"server"`
	Redis      Redis      `mapstructure:"redis"`
	DB         DB         `mapstructure:"db"`
	Caching    Caching    `mapstructure:"caching"`
	Service    Service    `mapstructure:"service"`
	AuthConfig AuthConfig `mapstructure:"auth"`
	Aud        string     `mapstructure:"aud"`
}

type Server struct {
	Host     string `mapstructure:"host"`
	GrpcPort string `mapstructure:"grpc_port"`
	HttpPort string `mapstructure:"http_port"`
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

type AuthConfig struct {
	Aud           string `mapstructure:"aud"`
	ExternalEnvoy struct {
		Header string `mapstructure:"header"`
		Scheme string `mapstructure:"scheme"`
	} `mapstructure:"external_envoy"`
	UserJWT struct {
		Type types.AuthUserType `mapstructure:"type"`
		RSA  struct {
			PrivateKeyPath string `mapstructure:"private_key_path"`
			PublicKeyPath  string `mapstructure:"public_key_path"`
		} `mapstructure:"rsa"`
		HMAC struct {
			Secret string `mapstructure:"secret"`
		} `mapstructure:"hmac"`
	} `mapstructure:"user_jwt"`
	Oauth struct {
		PublicKeyPath  string `mapstructure:"public_key_path"`
		PrivateKeyPath string `mapstructure:"private_key_path"`
		Clients        []struct {
			Name          string            `mapstructure:"name"`
			Type          types.AuthM2MType `mapstructure:"type"`
			PublicKey     string            `mapstructure:"public_key"`
			SecretKey     string            `mapstructure:"secret_key"`
			AllowUserAuth bool              `mapstructure:"allow_user_authenticate"`
		} `mapstructure:"clients"`
	} `mapstructure:"oauth"`
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
	v.SetConfigFile("config/config.yaml")

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
