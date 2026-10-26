package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/auth/v1"
)

func (s *GrpcServer) ListPermissions(ctx context.Context, req *authPb.ListPermissionsRequest) (*authPb.ListPermissionsResponse, error) {
	var keyword string
	var domains []string
	if filter := req.GetFilter(); filter != nil {
		keyword = filter.GetKeyword()
		domains = filter.GetDomains()
	}

	var pageNum, pageSize int64
	if pag := req.GetPagination(); pag != nil {
		pageNum = pag.GetPageNumber()
		pageSize = pag.GetPageSize()
	}

	res, err := s.svc.GetAuthService().ListPermissions(ctx, auth.ListPermissionsInput{
		Keyword:    keyword,
		Domains:    domains,
		PageNumber: pageNum,
		PageLimit:  pageSize,
	})
	if err != nil {
		return nil, err
	}

	return &authPb.ListPermissionsResponse{
		Permissions: ToPermissionsProto(res),
	}, nil
}
