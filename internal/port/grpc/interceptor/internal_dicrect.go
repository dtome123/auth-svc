package interceptor

import (
	"context"
	"fmt"
	"time"

	"auth-svc/config"
	"auth-svc/internal/types"

	"github.com/dtome123/auth-sdk/api/go/auth/v1"
	"github.com/dtome123/auth-sdk/client"
	"github.com/dtome123/auth-sdk/jwtutils"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// methodAuthRules defines per-method authentication rules for internal services.
var methodAuthRules = map[string]struct {
	RequireAssertion       bool // Whether this method requires internal client assertion
	RequireUserAuthSupport bool // Whether the client must support user authentication
}{
	auth.TokenService_Sign_FullMethodName: {
		RequireAssertion:       true,
		RequireUserAuthSupport: true,
	},
	auth.TokenService_Refresh_FullMethodName: {
		RequireAssertion:       true,
		RequireUserAuthSupport: true,
	},
	auth.AuthorizationService_MigratePermissions_FullMethodName: {
		RequireAssertion:       true,
		RequireUserAuthSupport: false,
	},
}

// internalDirectInterceptor handles internal JWT validation for gRPC services
type internalDirectInterceptor struct {
	clients       map[string]types.ClientEntry
	audience      string
	replayChecker jwtutils.ReplayChecker
}

// NewInternalDirectInterceptor constructs the interceptor and prebuilds JWT verifiers
func NewInternalDirectInterceptor(cfg config.AuthConfig, clients map[string]types.ClientEntry) *internalDirectInterceptor {

	var replayChecker jwtutils.ReplayChecker
	if cfg.EnableReplayCheck {
		replayChecker = jwtutils.NewMemoryReplayChecker(5 * time.Minute)
	}

	return &internalDirectInterceptor{
		clients:       clients,
		audience:      cfg.Aud,
		replayChecker: replayChecker,
	}
}

// UnaryInterceptor returns a gRPC interceptor that verifies client assertions
func (i *internalDirectInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		rule, protected := methodAuthRules[info.FullMethod]
		if !protected || !rule.RequireAssertion {
			return handler(ctx, req) // Not protected; allow through
		}

		clientID := client.ClientIDFromContext(ctx)
		entry, ok := i.clients[clientID]
		if !ok {
			return nil, fmt.Errorf("client %q is not whitelisted", clientID)
		}

		if rule.RequireUserAuthSupport && !entry.AllowUserAuthenticate {
			return nil, fmt.Errorf("client %q is not allowed to authenticate users", clientID)
		}

		token := client.ClientAssertionFromContext(ctx)
		if token == "" {
			return nil, errors.New("missing client assertion")
		}

		claims, err := entry.Verifier.Verify(token)
		if err != nil {
			return nil, errors.Wrap(err, "JWT verification failed")
		}

		oauthClaims := jwtutils.NewOauthClaims(claims)

		if err := oauthClaims.Validate(
			jwtutils.WithExpectedAudience(i.audience),
			jwtutils.WithReplayChecker(i.replayChecker),
		); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
