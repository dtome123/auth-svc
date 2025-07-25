// auth_interceptor.go
package interceptor

import (
	"auth-svc/config"
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	client "github.com/dtome123/auth-sdk/client"
	"github.com/dtome123/auth-sdk/constants"
	"github.com/dtome123/auth-sdk/jwtutils"
)

type clientAssertionInterceptor struct {
	whitelist     map[string]string
	aud           string
	ignoreMethods []string
}

func NewClientAssertionInterceptor(cfg config.AuthConfig) *clientAssertionInterceptor {

	whitelist := make(map[string]string)
	for _, svc := range cfg.M2M.Whitelist.InternalServices {
		whitelist[svc.Name] = svc.PublicKey
	}

	return &clientAssertionInterceptor{
		whitelist: whitelist,
		aud:       cfg.Aud,
		ignoreMethods: []string{
			"/envoy.service.auth.v3.Authorization/Check",
		},
	}
}

func (a *clientAssertionInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		method := info.FullMethod
		for _, m := range a.ignoreMethods {
			if method == m {
				return handler(ctx, req)
			}
		}

		clientID := client.ClientIDFromContext(ctx)
		if a.whitelist[clientID] == "" {
			return nil, fmt.Errorf("client %s is not whitelisted", clientID)
		}

		pubKey := a.whitelist[clientID]

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata")
		}

		values := md[constants.ClientAssertionKey]
		if len(values) == 0 {
			return nil, errors.New("missing client assertion")
		}

		tokenString := values[0]
		verifier, err := jwtutils.NewRS256VerifierFromString(pubKey)
		if err != nil {
			return nil, err
		}
		claims, err := verifier.Verify(tokenString)
		if err != nil {
			return nil, errors.New("invalid client assertion")
		}

		iss := claims.Get("iss").AsString()
		if iss == "" {
			return nil, errors.New("missing iss claim")
		}

		if aud := claims.Get("aud").AsString(); aud == "" || aud != a.aud {
			return nil, errors.New("invalid audience in assertion")
		}

		return handler(ctx, req)
	}
}
