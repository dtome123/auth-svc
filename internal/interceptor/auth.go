// auth_interceptor.go
package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	client "github.com/dtome123/auth-sdk/client"
	middleware "github.com/dtome123/auth-sdk/middlewares"
)

type ClientAssertionInterceptor struct {
	// Map client svc -> public key
	Whitelist map[string]string
	Aud       string
}

func NewClientAssertionInterceptor(whitelist map[string]string) (*ClientAssertionInterceptor, error) {

	return &ClientAssertionInterceptor{
		Whitelist: whitelist,
	}, nil
}

func (a *ClientAssertionInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		clientID := client.ClientIDFromContext(ctx)
		if a.Whitelist[clientID] == "" {
			return nil, fmt.Errorf("client %s is not whitelisted", clientID)
		}

		pubKey := a.Whitelist[clientID]
		middleware.ClientAssertionRSAUnaryInterceptor(a.Aud, pubKey)

		return handler(ctx, req)
	}
}
