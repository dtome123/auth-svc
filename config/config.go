package config

import (
	"auth-svc/internal/types"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server     Server     `mapstructure:"server"`
	Redis      Redis      `mapstructure:"redis"`
	DB         DB         `mapstructure:"db"`
	Caching    Caching    `mapstructure:"caching"`
	AuthConfig AuthConfig `mapstructure:"auth"`
}

type Server struct {
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
	Type   string `mapstructure:"type"`
	TTL    string `mapstructure:"ttl"`
}

type DB struct {
	Mongo struct {
		DSN      string `mapstructure:"dsn"`
		Database string `mapstructure:"database"`
	} `mapstructure:"mongo"`
}

type UserTokenConfig struct {
	Type            types.AuthUserType `mapstructure:"type"`
	AccessTokenTTL  string             `mapstructure:"access_token_ttl"`
	RefreshTokenTTL string             `mapstructure:"refresh_token_ttl"`
	RSA             struct {
		PrivateKeyPath string `mapstructure:"private_key_path"`
		PublicKeyPath  string `mapstructure:"public_key_path"`
	} `mapstructure:"rsa"`
	HMAC struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"hmac"`
}

func (u *UserTokenConfig) IsEmpty() bool {
	if u == nil {
		return true
	}
	return u.Type == ""
}

type ClientConfig struct {
	Name            string            `mapstructure:"name"`
	Type            types.AuthM2MType `mapstructure:"type"`
	PublicKey       string            `mapstructure:"public_key"`
	SecretKey       string            `mapstructure:"secret_key"`
	AllowUserAuth   bool              `mapstructure:"allow_user_authenticate"`
	AllowAudiences  []string          `mapstructure:"allow_audiences"`
	UserTokenConfig *UserTokenConfig  `mapstructure:"user_token_config"`
}

type AuthConfig struct {
	SystemAud     string `mapstructure:"system_aud"`
	ExternalEnvoy struct {
		Header string `mapstructure:"header"`
		Scheme string `mapstructure:"scheme"`
	} `mapstructure:"external_envoy"`
	EnableReplayCheck bool `mapstructure:"enable_replay_check"`
	Oauth             struct {
		Aud  string            `mapstructure:"aud"`
		Type types.AuthM2MType `mapstructure:"type"`
		HMAC struct {
			Secret string `mapstructure:"secret"`
		} `mapstructure:"hmac"`
		RSA struct {
			PrivateKeyPath string `mapstructure:"private_key_path"`
			PublicKeyPath  string `mapstructure:"public_key_path"`
		} `mapstructure:"rsa"`
		AccessTokenTTL  string          `mapstructure:"access_token_ttl"`
		UserTokenConfig UserTokenConfig `mapstructure:"user_token_config"`
		Clients         []ClientConfig  `mapstructure:"clients"`
	} `mapstructure:"oauth"`
}

func (c *Config) Validate() error {
	if c.AuthConfig.Oauth.UserTokenConfig.IsEmpty() {
		return fmt.Errorf("user_token_config in oauth cannot be empty")
	}

	for i := range c.AuthConfig.Oauth.Clients {
		client := &c.AuthConfig.Oauth.Clients[i]
		if client.AllowUserAuth {
			if client.UserTokenConfig == nil {
				client.UserTokenConfig = &c.AuthConfig.Oauth.UserTokenConfig
			} else {
				if client.UserTokenConfig.Type == "" {
					client.UserTokenConfig.Type = c.AuthConfig.Oauth.UserTokenConfig.Type
					client.UserTokenConfig.RSA = c.AuthConfig.Oauth.UserTokenConfig.RSA
					client.UserTokenConfig.HMAC = c.AuthConfig.Oauth.UserTokenConfig.HMAC
				}
				if client.UserTokenConfig.AccessTokenTTL == "" {
					client.UserTokenConfig.AccessTokenTTL = c.AuthConfig.Oauth.UserTokenConfig.AccessTokenTTL
				}
				if client.UserTokenConfig.RefreshTokenTTL == "" {
					client.UserTokenConfig.RefreshTokenTTL = c.AuthConfig.Oauth.UserTokenConfig.RefreshTokenTTL
				}
			}
		}
	}

	return nil
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

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
