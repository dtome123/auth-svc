package main

import (
	"auth-svc/config"
	"auth-svc/internal/domain"
	"auth-svc/internal/port"
	"auth-svc/internal/types"
	"log"

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
		var signer jwt.Signer
		var err error

		switch svc.Type {
		case types.AuthM2MTypeRSA:
			verifier, err = jwt.NewRS256VerifierFromString(svc.PublicKey)
			if err != nil {
				log.Fatalf("Invalid RSA key for client %q: %v", svc.Name, err)

			}
			signer, err = jwt.NewRS256SignerFromString(svc.SecretKey)
			if err != nil {
				log.Fatalf("Invalid RSA key for client %q: %v", svc.Name, err)
			}
		case types.AuthM2MTypeHMAC:
			verifier = jwt.NewHMACVerifier([]byte(svc.SecretKey))
			signer = jwt.NewHMACSigner([]byte(svc.SecretKey))
		default:
			log.Fatalf("Unsupported auth type %q for client %q", svc.Type, svc.Name)
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
			Signer:                signer,
			Verifier:              verifier,
		}
	}

	return clients
}
