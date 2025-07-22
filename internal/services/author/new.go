package author

import (
	authorDb "auth-svc/internal/adapters/mongodb/author"
	authorCache "auth-svc/internal/adapters/redis/author"
)

type AuthorizationService struct {
	Repo  *authorDb.AuthorizationRepository
	Cache *authorCache.AuthorizationCacheRepository
}

func NewAuthorizationService(repo *authorDb.AuthorizationRepository, cache *authorCache.AuthorizationCacheRepository) *AuthorizationService {
	return &AuthorizationService{
		Repo:  repo,
		Cache: cache,
	}
}
