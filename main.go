package main

import (
	"auth-svc/config"
	"auth-svc/internal/domain"
	"auth-svc/internal/port"
	"auth-svc/internal/types"
	"fmt"

	"github.com/dtome123/auth-sdk/jwt"
	"github.com/dtome123/auth-sdk/session"
	mongodb "github.com/dtome123/go-mongo-generic"
	"github.com/redis/go-redis/v9"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := mongodb.NewDatabase(
		mongodb.WithDatabase(cfg.DB.Mongo.Database),
		mongodb.WithSingleURL(cfg.DB.Mongo.DSN),
	)
	if err != nil {
		panic(err)
	}

	var redisClient *redis.Client
	var sessionResource session.Repository[domain.Session]
	if cfg.Redis.Host != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
			Password: cfg.Redis.PWD,
			DB:       cfg.Redis.Database,
		})

		sessionResource = session.New[domain.Session](redisClient)
	}

	clientVerifiers := newClientVerifiers(cfg.AuthConfig)

	server := port.NewServer(cfg, db, redisClient, clientVerifiers, sessionResource)
	server.Run()
}

func newClientVerifiers(authConfig config.AuthConfig) map[string]types.ClientEntry {
	clients := make(map[string]types.ClientEntry)

	for _, svc := range authConfig.Oauth.Clients {
		var verifier jwt.Verifier
		var err error

		switch svc.Type {
		case types.AuthM2MTypeRSA:
			verifier, err = jwt.NewRS256VerifierFromString(svc.PublicKey)
			if err != nil {
				panic(fmt.Sprintf("Invalid RSA key for client %q: %v", svc.Name, err))
			}
		case types.AuthM2MTypeHMAC:
			verifier = jwt.NewHMACVerifier([]byte(svc.SecretKey))
		default:
			panic(fmt.Sprintf("Unsupported auth type %q for client %q", svc.Type, svc.Name))
		}

		var userSigner jwt.Signer
		var userVerifier jwt.Verifier
		if svc.UserTokenConfig != nil && !svc.UserTokenConfig.IsEmpty() {
			switch svc.UserTokenConfig.Type {
			case types.AuthUserTypeHMAC:
				userSigner = jwt.NewHMACSigner([]byte(svc.UserTokenConfig.HMAC.Secret))
				userVerifier = jwt.NewHMACVerifier([]byte(svc.UserTokenConfig.HMAC.Secret))
			case types.AuthUserTypeRSA:
				userSigner, err = jwt.NewRS256SignerFromPath(svc.UserTokenConfig.RSA.PrivateKeyPath)
				if err != nil {
					panic(fmt.Sprintf("Invalid RSA private key for client %q: %v", svc.Name, err))
				}
				userVerifier, err = jwt.NewRS256VerifierFromPath(svc.UserTokenConfig.RSA.PublicKeyPath)
				if err != nil {
					panic(fmt.Sprintf("Invalid RSA public key for client %q: %v", svc.Name, err))
				}
			}
		}

		var accessTokenTTL string
		var refreshTokenTTL string
		if svc.UserTokenConfig != nil && !svc.UserTokenConfig.IsEmpty() {
			accessTokenTTL = svc.UserTokenConfig.AccessTokenTTL
			refreshTokenTTL = svc.UserTokenConfig.RefreshTokenTTL
		}

		allowAudiences := make(map[string]bool)
		for _, audience := range svc.AllowAudiences {
			allowAudiences[audience] = true
		}

		clients[svc.Name] = types.ClientEntry{
			Name:                  svc.Name,
			Type:                  svc.Type,
			AllowUserAuthenticate: svc.AllowUserAuth,
			AllowAudiences:        allowAudiences,
			AccessTokenTTL:        accessTokenTTL,
			RefreshTokenTTL:       refreshTokenTTL,
			Verifier:              verifier,
			UserSigner:            userSigner,
			UserVerifier:          userVerifier,
		}
	}

	return clients
}
