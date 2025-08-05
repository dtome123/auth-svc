package main

import (
	"auth-svc/config"
	"auth-svc/internal/port"
	"auth-svc/internal/types"
	"fmt"

	"github.com/dtome123/auth-sdk/jwtutils"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.PWD,
		DB:       cfg.Redis.Database,
	})

	clientVerifiers := newClientVerifiers(cfg.AuthConfig)

	server := port.NewServer(cfg, db, redisClient, clientVerifiers)
	server.Run()
}

func newClientVerifiers(authConfig config.AuthConfig) map[string]types.ClientEntry {
	clients := make(map[string]types.ClientEntry)

	for _, svc := range authConfig.Oauth.Clients {
		var verifier jwtutils.Verifier
		var err error

		switch svc.Type {
		case types.AuthM2MTypeRSA:
			verifier, err = jwtutils.NewRS256VerifierFromString(svc.PublicKey)
			if err != nil {
				panic(fmt.Sprintf("Invalid RSA key for client %q: %v", svc.Name, err))
			}
		case types.AuthM2MTypeHMAC:
			verifier = jwtutils.NewHMACVerifier([]byte(svc.SecretKey))
		default:
			panic(fmt.Sprintf("Unsupported auth type %q for client %q", svc.Type, svc.Name))
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
			Verifier:              verifier,
		}
	}

	return clients
}
