package grpc

import (
	"auth-svc/internal/services/auth"
	"auth-svc/internal/utils"
	"context"

	authPb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/genproto/googleapis/rpc/status"
)

func (s *GrpcServer) Check(ctx context.Context, req *authPb.CheckRequest) (*authPb.CheckResponse, error) {

	path := req.Attributes.Request.Http.Path
	authHeader := req.Attributes.Request.Http.GetHeaders()

	parse := utils.ExtractExternalToken(authHeader, utils.AuthConfig{
		Header: s.cfg.AuthConfig.ExternalEnvoy.Header,
		Scheme: s.cfg.AuthConfig.ExternalEnvoy.Scheme,
	})

	result, err := s.svc.GetAuthService().Check(ctx, auth.CheckInput{
		AccessToken: parse.Token,
		FullMethod:  path,
	})

	if err != nil {
		return nil, err
	}

	if !result.Allowed {
		return &authPb.CheckResponse{
			Status: &status.Status{
				Code:    result.StatusCode,
				Message: result.Message,
			},
			HttpResponse: &authPb.CheckResponse_DeniedResponse{
				DeniedResponse: &authPb.DeniedHttpResponse{
					Status: &typev3.HttpStatus{
						Code: typev3.StatusCode(result.StatusCode),
					},
					Body: result.Message,
				},
			},
		}, nil
	}

	return &authPb.CheckResponse{
		Status: &status.Status{
			Code:    result.StatusCode,
			Message: result.Message,
		},
		HttpResponse: &authPb.CheckResponse_OkResponse{
			OkResponse: &authPb.OkHttpResponse{},
		},
	}, nil
}
