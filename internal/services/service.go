package services

import (
	"auth-svc/config"
	"auth-svc/internal/adapters/mongodb/authen"
	authorDb "auth-svc/internal/adapters/mongodb/author"
	"auth-svc/internal/adapters/redis/author"
	"auth-svc/internal/services/auth"
	"auth-svc/internal/types"
	"time"

	"github.com/dtome123/auth-sdk/jwtutils"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	cfg         *config.Config
	clients     map[string]types.ClientEntry
	authService *auth.AuthorizationService
}

func (s *Service) GetAuthService() *auth.AuthorizationService {
	return s.authService
}

func NewService(
	cfg *config.Config,
	db *mongo.Database,
	redisClient *redis.Client,
	clientVerifiers map[string]types.ClientEntry,
	serverSigner jwtutils.Signer,
	serverVerifier jwtutils.Verifier,
) *Service {

	authorRepo := authorDb.NewAuthorizationRepository(db)
	authenRepo := authen.NewAuthenticationRepository(db)
	authorCacheRepo := author.NewAuthorizationCacheRepository(redisClient, cfg.Caching.TTL)

	return &Service{
		cfg:     cfg,
		clients: clientVerifiers,
		authService: auth.NewAuthorizationService(
			cfg,
			serverSigner,
			serverVerifier,
			authorRepo,
			authenRepo,
			authorCacheRepo,
			jwtutils.NewMemoryReplayChecker(5*time.Minute),
			clientVerifiers,
		),
	}
}
