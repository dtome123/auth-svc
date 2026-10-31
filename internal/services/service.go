package services

import (
	"auth-svc/config"
	"auth-svc/internal/adapters/cache"
	authorDb "auth-svc/internal/adapters/mongodb/author"
	"auth-svc/internal/domain"
	"auth-svc/internal/services/auth"
	"auth-svc/internal/types"
	"time"

	"github.com/dtome123/auth-sdk/assertion"
	"github.com/dtome123/auth-sdk/jwt"
	"github.com/dtome123/auth-sdk/session"
	mongodb "github.com/dtome123/go-mongo-generic"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	cfg             *config.Config
	clients         map[string]types.ClientEntry
	authService     *auth.AuthorizationService
	replayChecker   assertion.ReplayChecker
	sessionResource session.Repository[domain.Session]
}

func (s *Service) GetAuthService() *auth.AuthorizationService {
	return s.authService
}

func (s *Service) GetReplayChecker() assertion.ReplayChecker {
	return s.replayChecker
}

func NewService(
	cfg *config.Config,
	db *mongodb.Database,
	redisClient *redis.Client,
	clientVerifiers map[string]types.ClientEntry,
	serverSigner jwt.Signer,
	serverVerifier jwt.Verifier,
	userSigner jwt.Signer,
	userVerifier jwt.Verifier,
	sessionResource session.Repository[domain.Session],
) *Service {

	authorRepo := authorDb.NewAuthorizationRepository(db)

	var authorCacheRepo cache.AuthorizationCacheRepository
	if !cfg.Caching.Enable || cfg.Caching.Type == "none" {
		authorCacheRepo = cache.NewNoOpAuthorizationCacheRepository()
	} else if cfg.Caching.Type == "memory" {
		authorCacheRepo = cache.NewMemoryAuthorizationCacheRepository(cfg.Caching.TTL)
	} else {
		if redisClient != nil {
			authorCacheRepo = cache.NewRedisAuthorizationCacheRepository(redisClient, cfg.Caching.TTL)
		} else {
			authorCacheRepo = cache.NewNoOpAuthorizationCacheRepository()
		}
	}

	var replayChecker assertion.ReplayChecker
	if cfg.AuthConfig.EnableReplayCheck {
		if cfg.Caching.Enable && cfg.Caching.Type == "redis" && redisClient != nil {
			replayChecker = assertion.NewRedisReplayChecker(redisClient)
		} else {
			replayChecker = assertion.NewMemoryReplayChecker(5 * time.Minute)
		}
	}

	return &Service{
		cfg:           cfg,
		clients:       clientVerifiers,
		replayChecker: replayChecker,
		authService: auth.NewAuthorizationService(
			cfg,
			serverSigner,
			serverVerifier,
			userSigner,
			userVerifier,
			authorRepo,
			sessionResource,
			authorCacheRepo,
			replayChecker,
			clientVerifiers,
			redisClient,
		),
	}
}
