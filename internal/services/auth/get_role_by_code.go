package auth

import (
	"auth-svc/internal/domain"
	"context"

	"github.com/pkg/errors"
)

func (svc *AuthorizationService) GetRoleByCode(ctx context.Context, code string, scope map[string][]string) (*domain.Role, error) {
	role, err := svc.authorizationRepo.GetRoleByCode(ctx, code, scope)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get role by code")
	}
	if role == nil {
		return nil, nil
	}

	return role, nil
}
