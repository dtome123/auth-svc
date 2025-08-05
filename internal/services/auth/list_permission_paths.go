package auth

import (
	"auth-svc/internal/adapters/mongodb/author"
	"auth-svc/internal/models"
	"context"

	"github.com/pkg/errors"
)

type ListPermissionPathsInput struct {
	Domains []string
	Keyword string
}

func (s *AuthorizationService) ListPermissionPaths(ctx context.Context, input ListPermissionPathsInput) ([]*models.PermissionPath, error) {

	paths, err := s.authorizationRepo.ListPermissionPaths(ctx, author.ListPermissionPathsInput{
		Domains: input.Domains,
		Keyword: input.Keyword,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to list permission paths")
	}

	return paths, nil
}
