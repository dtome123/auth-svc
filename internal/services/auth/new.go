package auth

import (
	"auth-svc/config"
	"auth-svc/internal/adapters/cache"
	"auth-svc/internal/adapters/mongodb/author"
	"auth-svc/internal/domain"
	"auth-svc/internal/types"
	casbinPkg "auth-svc/tools/casbin"
	"context"
	"log"
	"time"

	"github.com/dtome123/auth-sdk/assertion"
	"github.com/dtome123/auth-sdk/jwt"
	"github.com/dtome123/auth-sdk/session"
	"github.com/redis/go-redis/v9"
)

type AuthorizationService struct {
	cfg               *config.Config
	clients           map[string]types.ClientEntry
	serverSigner      jwt.Signer
	serverVerifier    jwt.Verifier
	userSigner        jwt.Signer
	userVerifier      jwt.Verifier
	authorizationRepo author.AuthorizationRepository
	sessionResource   session.Repository[domain.Session]
	authorCache       cache.AuthorizationCacheRepository
	oauthRelayChecker assertion.ReplayChecker
	enforcer          casbinPkg.Enforcer
}

func NewAuthorizationService(
	cfg *config.Config,
	serverSigner jwt.Signer,
	serverVerifier jwt.Verifier,
	userSigner jwt.Signer,
	userVerifier jwt.Verifier,
	authorizationRepo author.AuthorizationRepository,
	sessionResource session.Repository[domain.Session],
	authorCache cache.AuthorizationCacheRepository,
	oauthRelayChecker assertion.ReplayChecker,
	clientVerifiers map[string]types.ClientEntry,
	redisClient *redis.Client,
) *AuthorizationService {

	// Initialize Casbin enforcer
	casbinCfg := casbinPkg.Config{
		SyncChannel: "auth:casbin:sync",
	}

	var enforcerOpts []casbinPkg.Option
	if redisClient != nil {
		enforcerOpts = append(enforcerOpts, casbinPkg.WithRedis(redisClient))
	}

	enforcer, err := casbinPkg.New(casbinCfg, enforcerOpts...)
	if err != nil {
		panic(err)
	}

	svc := &AuthorizationService{
		cfg:               cfg,
		clients:           clientVerifiers,
		serverSigner:      serverSigner,
		serverVerifier:    serverVerifier,
		userSigner:        userSigner,
		userVerifier:      userVerifier,
		sessionResource:   sessionResource,
		authorizationRepo: authorizationRepo,
		authorCache:       authorCache,
		oauthRelayChecker: oauthRelayChecker,
		enforcer:          enforcer,
	}

	// Register reload callback and start subscriber goroutine
	enforcer.OnSyncEvent(svc.LoadAllPolicies)

	// Warm-up policies from MongoDB on startup
	if authorizationRepo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := svc.LoadAllPolicies(ctx); err != nil {
			log.Printf("⚠️ Warning: failed to load casbin policies on startup: %v", err)
		} else {
			log.Printf("✅ Casbin RBAC policies loaded successfully into memory")
		}
	}

	return svc
}
