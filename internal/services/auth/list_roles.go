package auth

import (
	"auth-svc/internal/domain"
	"context"

	"github.com/pkg/errors"
)

type ListRoleInput = domain.RoleCriteria

func (s *AuthorizationService) ListRole(ctx context.Context, input ListRoleInput) ([]*domain.Role, error) {
	roles, err := s.authorizationRepo.ListRoles(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list roles")
	}

	return roles, nil
}
