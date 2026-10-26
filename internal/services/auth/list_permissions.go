package auth

import (
	"auth-svc/internal/domain"
	"context"

	"github.com/pkg/errors"
)

type ListPermissionsInput = domain.ListPermissionsInput

func (s *AuthorizationService) ListPermissions(ctx context.Context, input ListPermissionsInput) ([]*domain.Permission, error) {

	permissions, err := s.authorizationRepo.ListPermissions(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list permissions")
	}

	return permissions, nil
}
