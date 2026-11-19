package auth

import (
	"auth-svc/internal/domain"
	"context"

	"github.com/pkg/errors"
)

func (svc *AuthorizationService) GetRole(ctx context.Context, roleID string) (*domain.Role, error) {
	role, err := svc.authorizationRepo.GetRole(ctx, roleID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get role")
	}
	if role == nil {
		return nil, nil
	}

	return role, nil
}
