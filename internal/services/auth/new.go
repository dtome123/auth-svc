package auth

import (
	"auth-svc/config"
	"auth-svc/internal/adapters/mongodb/authen"
	authorDb "auth-svc/internal/adapters/mongodb/author"
	authorCache "auth-svc/internal/adapters/redis/author"
	"auth-svc/internal/types"

	"github.com/dtome123/auth-sdk/jwtutils"
)

type AuthorizationService struct {
	cfg                *config.Config
	clients            map[string]types.ClientEntry
	serverSigner       jwtutils.Signer
	serverVerifier     jwtutils.Verifier
	authorizationRepo  *authorDb.AuthorizationRepository
	authenticationRepo *authen.AuthenticationRepository
	authorCache        *authorCache.AuthorizationCacheRepository
	oauthRelayChecker  jwtutils.ReplayChecker
}

func NewAuthorizationService(
	cfg *config.Config,
	serverSigner jwtutils.Signer,
	serverVerifier jwtutils.Verifier,
	authorizationRepo *authorDb.AuthorizationRepository,
	authenticationRepo *authen.AuthenticationRepository,
	authorCache *authorCache.AuthorizationCacheRepository,
	oauthRelayChecker jwtutils.ReplayChecker,
	clientVerifiers map[string]types.ClientEntry,
) *AuthorizationService {

	return &AuthorizationService{
		cfg:                cfg,
		clients:            clientVerifiers,
		serverSigner:       serverSigner,
		serverVerifier:     serverVerifier,
		authenticationRepo: authenticationRepo,
		authorizationRepo:  authorizationRepo,
		authorCache:        authorCache,
		oauthRelayChecker:  oauthRelayChecker,
	}
}
