package services

import (
	"auth-svc/config"
	"auth-svc/internal/adapters/mongodb/authen"
	authorDb "auth-svc/internal/adapters/mongodb/author"
	"auth-svc/internal/adapters/redis/author"
	"auth-svc/internal/services/auth"
	"auth-svc/internal/types"

	"github.com/dtome123/auth-sdk/jwtutils"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	cfg         *config.Config
	authService *auth.AuthorizationService
}

func (s *Service) GetAuthService() *auth.AuthorizationService {
	return s.authService
}

func NewService(cfg *config.Config, db *mongo.Database, redisClient *redis.Client) *Service {

	authorRepo := authorDb.NewAuthorizationRepository(db)
	authenRepo := authen.NewAuthenticationRepository(db)
	var authorCacheRepo *author.AuthorizationCacheRepository

	if cfg.Caching.Enable {
		authorCacheRepo = author.NewAuthorizationCacheRepository(redisClient, cfg.Caching.TTL)
	}
	var serverSigner jwtutils.Signer
	var serverVerifier jwtutils.Verifier

	switch cfg.AuthConfig.UserJWT.Type {
	case types.AuthUserTypeHMAC:
		serverSigner = jwtutils.NewHMACSigner([]byte(cfg.AuthConfig.UserJWT.HMAC.Secret))
		serverVerifier = jwtutils.NewHMACVerifier([]byte(cfg.AuthConfig.UserJWT.HMAC.Secret))
	case types.AuthUserTypeRSA:
		var err error
		serverSigner, err = jwtutils.NewRS256SignerFromPath(cfg.AuthConfig.UserJWT.RSA.PrivateKeyPath)
		if err != nil {
			panic(err)
		}

		serverVerifier, err = jwtutils.NewRS256VerifierFromPath(cfg.AuthConfig.UserJWT.RSA.PublicKeyPath)
		if err != nil {
			panic(err)
		}
	}

	return &Service{
		cfg:         cfg,
		authService: auth.NewAuthorizationService(cfg, serverSigner, serverVerifier, authorRepo, authenRepo, authorCacheRepo),
	}
}
