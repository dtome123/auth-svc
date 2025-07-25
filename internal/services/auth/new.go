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
	clients            map[string]types.AuthClientEntry
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

	clients := make(map[string]types.AuthClientEntry)
	for _, client := range cfg.AuthConfig.Oauth.Clients {
		clients[client.Name] = types.AuthClientEntry{
			Name:      client.Name,
			Type:      client.Type,
			PublicKey: client.PublicKey,
			SecretKey: client.SecretKey,
		}
	}

	return &AuthorizationService{
		cfg:                cfg,
		clients:            clients,
		serverSigner:       serverSigner,
		serverVerifier:     serverVerifier,
		authenticationRepo: authenticationRepo,
		authorizationRepo:  authorizationRepo,
		authorCache:        authorCache,
	}
}
