package auth

import (
	"auth-svc/config"
	"auth-svc/internal/adapters/mongodb/authen"
	authorDb "auth-svc/internal/adapters/mongodb/author"
	authorCache "auth-svc/internal/adapters/redis/author"

	"github.com/dtome123/auth-sdk/jwtutils"
)

type AuthorizationService struct {
	cfg                *config.Config
	serverSigner       jwtutils.Signer
	serverVerifier     jwtutils.Verifier
	authorizationRepo  *authorDb.AuthorizationRepository
	authenticationRepo *authen.AuthenticationRepository
	authorCache        *authorCache.AuthorizationCacheRepository
}

func NewAuthorizationService(
	cfg *config.Config,
	serverSigner jwtutils.Signer,
	serverVerifier jwtutils.Verifier,
	authorizationRepo *authorDb.AuthorizationRepository,
	authenticationRepo *authen.AuthenticationRepository,
	authorCache *authorCache.AuthorizationCacheRepository,
) *AuthorizationService {
	return &AuthorizationService{
		cfg:                cfg,
		serverSigner:       serverSigner,
		serverVerifier:     serverVerifier,
		authenticationRepo: authenticationRepo,
		authorizationRepo:  authorizationRepo,
		authorCache:        authorCache,
	}
}
