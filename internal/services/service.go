package services

import (
	"auth-svc/config"
	"auth-svc/internal/adapters/mongodb/authen"
	authorDb "auth-svc/internal/adapters/mongodb/author"
	"auth-svc/internal/adapters/redis/author"
	"auth-svc/internal/services/auth"
	"crypto/sha256"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	jwtCrypto "github.com/dtome123/auth-sdk/jwtutils/crypto"
)

type Service struct {
	cfg         *config.Config
	authService *auth.AuthorizationService
	jwtCrypto   jwtCrypto.JWTCrypto
}

func NewService(cfg *config.Config, db *mongo.Database, redisClient *redis.Client) *Service {

	authorRepo := authorDb.NewAuthorizationRepository(db)
	authenRepo := authen.NewAuthenticationRepository(db)
	var authorCacheRepo *author.AuthorizationCacheRepository

	if cfg.Caching.Enable {
		authorCacheRepo = author.NewAuthorizationCacheRepository(redisClient, cfg.Caching.TTL)
	}

	return &Service{
		cfg: cfg,
		jwtCrypto: jwtCrypto.NewRsaOEAPJWTCrypto(jwtCrypto.RsaOEAPJWTConfig{
			PubPath:       cfg.Cert.PublicKeyPath,
			PrivPath:      cfg.Cert.PrivateKeyPath,
			OaepLabel:     []byte("auth-svc"),
			OaepHashNewFn: sha256.New,
		}),
		authService: auth.NewAuthorizationService(cfg, authorRepo, authenRepo, authorCacheRepo),
	}
}
