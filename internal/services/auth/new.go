package auth

import (
	"auth-svc/config"
	"auth-svc/internal/adapters/mongodb/authen"
	authorDb "auth-svc/internal/adapters/mongodb/author"
	authorCache "auth-svc/internal/adapters/redis/author"
)

type AuthorizationService struct {
	cfg                *config.Config
	authorizationRepo  *authorDb.AuthorizationRepository
	authenticationRepo *authen.AuthenticationRepository
	authorCache        *authorCache.AuthorizationCacheRepository
}

func NewAuthorizationService(
	cfg *config.Config,
	authorizationRepo *authorDb.AuthorizationRepository,
	authenticationRepo *authen.AuthenticationRepository,
	authorCache *authorCache.AuthorizationCacheRepository,
) *AuthorizationService {
	return &AuthorizationService{
		cfg:               cfg,
		authenticationRepo: authenticationRepo,
		authorizationRepo: authorizationRepo,
		authorCache:       authorCache,
	}
}
