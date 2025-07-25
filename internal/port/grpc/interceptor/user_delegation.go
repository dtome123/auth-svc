package interceptor

import (
	"context"
	"fmt"
	"time"

	"auth-svc/config"
	"auth-svc/internal/types"

	"github.com/dtome123/auth-sdk/api/go/auth/v1"
	"github.com/dtome123/auth-sdk/client"
	"github.com/dtome123/auth-sdk/constants"
	"github.com/dtome123/auth-sdk/jwtutils"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Define allowed routes for different scopes
var (
	userAuthScopeRoutes = map[string]struct{}{
		auth.AuthService_Sign_FullMethodName:    {},
		auth.AuthService_Refresh_FullMethodName: {},
	}
)

// UserAuthClientEntry represents a trusted client allowed to authenticate using client assertion (JWT)
type UserAuthClientEntry struct {
	Name                  string            `json:"name"`       // Service name
	Type                  types.AuthM2MType `json:"type"`       // Auth method: "rsa" or "hmac"
	PublicKey             string            `json:"public_key"` // RSA public key (for verification)
	SecretKey             string            `json:"secret_key"` // HMAC secret key
	AllowUserAuthenticate bool              `json:"allow_user_authenticate"`
}

// UserDelegationInterceptor is a gRPC interceptor that verifies client assertion tokens (JWTs)
type UserDelegationInterceptor struct {
	userAuthClient map[string]UserAuthClientEntry // Map of whitelisted client IDs
	audience       string                         // Expected "aud" claim in JWTs
	jtiCache       *cache.Cache                   // In-memory cache to prevent replay attacks
}

// NewUserDelegationInterceptor initializes a new interceptor with whitelist and audience settings
func NewUserDelegationInterceptor(cfg config.AuthConfig) *UserDelegationInterceptor {
	userAuthClient := make(map[string]UserAuthClientEntry)
	for _, svc := range cfg.Oauth.Clients {
		userAuthClient[svc.Name] = UserAuthClientEntry{
			Name:                  svc.Name,
			Type:                  svc.Type,
			PublicKey:             svc.PublicKey,
			SecretKey:             svc.SecretKey,
			AllowUserAuthenticate: svc.AllowUserAuth,
		}
	}

	return &UserDelegationInterceptor{
		userAuthClient: userAuthClient,
		audience:       cfg.Aud,
		jtiCache:       cache.New(5*time.Minute, 10*time.Minute),
	}
}

// UnaryInterceptor returns a gRPC UnaryServerInterceptor for validating client assertion JWTs
func (interceptor *UserDelegationInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// Step 1: Skip verification if the method is not protected
		if _, ok := userAuthScopeRoutes[info.FullMethod]; !ok {
			return handler(ctx, req)
		}

		// Step 2: Extract client ID from context and verify against whitelist
		clientID := client.ClientIDFromContext(ctx)
		issuer, ok := interceptor.userAuthClient[clientID]
		if !ok {
			return nil, fmt.Errorf("client %q is not whitelisted", clientID)
		}

		// Step 3: Check if client is allowed to authenticate
		if !issuer.AllowUserAuthenticate {
			return nil, fmt.Errorf("client %q is not allowed to authenticate", clientID)
		}

		// Step 4: Extract JWT from gRPC metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata")
		}
		assertionTokens := md[constants.ClientAssertionKey]
		if len(assertionTokens) == 0 || assertionTokens[0] == "" {
			return nil, errors.New("missing client assertion")
		}
		assertionJWT := assertionTokens[0]

		// Step 5: Create JWT verifier based on auth type
		var verifier jwtutils.Verifier
		var err error
		switch issuer.Type {
		case types.AuthM2MTypeRSA:
			verifier, err = jwtutils.NewRS256VerifierFromString(issuer.PublicKey)
			if err != nil {
				return nil, errors.Wrap(err, "invalid public key")
			}
		case types.AuthM2MTypeHMAC:
			verifier = jwtutils.NewHMACVerifier([]byte(issuer.SecretKey))
		default:
			return nil, errors.Errorf("unsupported auth type %q for client %q", issuer.Type, issuer.Name)
		}

		// Step 6: Verify JWT and extract claims
		claims, err := verifier.Verify(assertionJWT)
		if err != nil {
			return nil, errors.Wrap(err, "JWT verification failed")
		}

		// Step 7: Validate standard claims
		if claims.Get("iss").AsString() == "" {
			return nil, errors.New("missing 'iss' claim")
		}
		if aud := claims.Get("aud").AsString(); aud != interceptor.audience {
			return nil, errors.New("invalid or missing 'aud' claim")
		}
		exp := claims.Get("exp").AsInt64()
		if exp == 0 {
			return nil, errors.New("missing 'exp' claim")
		}
		now := time.Now().Unix()
		if exp < now {
			return nil, errors.New("JWT expired")
		}
		jti := claims.Get("jti").AsString()
		if jti == "" {
			return nil, errors.New("missing 'jti' claim")
		}

		// Step 8: Prevent token replay by checking JTI
		if _, expiry, exists := interceptor.jtiCache.GetWithExpiration(jti); exists && expiry.Unix() > now {
			return nil, errors.New("JWT replay detected (duplicate jti)")
		}

		// Step 9: Cache the JTI until token expiry
		ttl := time.Duration(exp-now) * time.Second
		interceptor.jtiCache.Set(jti, true, ttl)

		// Step 10: Proceed with actual handler
		return handler(ctx, req)
	}
}
