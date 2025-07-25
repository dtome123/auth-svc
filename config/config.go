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
	Client struct {
		Type types.AuthClient `mapstructure:"type"`
		RSA  struct {
			PrivateKeyPath string `mapstructure:"private_key_path"`
			PublicKeyPath  string `mapstructure:"public_key_path"`
		} `mapstructure:"rsa"`
		HMAC struct {
			Secret string `mapstructure:"secret"`
		} `mapstructure:"hmac"`
	} `mapstructure:"client"`

	M2M struct {
		EnableAssertion bool             `mapstructure:"enable_assertion"`
		Type            types.AuthClient `mapstructure:"type"`
		Whitelist       Whitelist        `mapstructure:"whitelist"`
	} `mapstructure:"m2m"`
}

type Service struct {
	Session struct {
		AccessTokenTTL  string `mapstructure:"access_token_ttl"`
		RefreshTokenTTL string `mapstructure:"refresh_token_ttl"`
	} `mapstructure:"session"`
}

type Whitelist struct {
	InternalServices []struct {
		Name      string `mapstructure:"name"`
		PublicKey string `mapstructure:"public_key"`
	} `mapstructure:"internal_services"`
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
