package auth

import (
	"auth-svc/internal/adapters/mongodb/author"
	"auth-svc/internal/models"
	"context"

	"github.com/pkg/errors"
)

type ListPermissionsInput struct {
	Domains []string
	Keyword string
}

func (s *AuthorizationService) ListPermissions(ctx context.Context, input ListPermissionsInput) ([]*models.Permission, error) {

	permissions, err := s.authorizationRepo.ListPermissions(ctx, author.ListPermissionsInput{
		Domains: input.Domains,
		Keyword: input.Keyword,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to list permission paths")
	}

	return permissions, nil
}
