package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
)

func (s *GrpcServer) Token(ctx context.Context, req *authPb.TokenRequest) (*authPb.TokenResponse, error) {

	res, err := s.svc.GetAuthService().Token(ctx, auth.TokenInput{
		GrantType:           req.GrantType,
		ClientAssertionType: req.ClientAssertionType,
		ClientAssertion:     req.ClientAssertion,
	})

	if err != nil {
		return nil, err
	}

	return &authPb.TokenResponse{
		AccessToken: res.AccessToken,
		TokenType:   res.TokenType,
		ExpiresIn:   res.ExpiresIn,
	}, nil
}
