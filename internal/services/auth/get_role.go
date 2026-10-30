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

	if len(role.ResourceActions) > 0 {
		permIDs := make([]string, 0, len(role.ResourceActions))
		for _, ra := range role.ResourceActions {
			if ra.PermissionID != "" {
				permIDs = append(permIDs, ra.PermissionID)
			}
		}
		if len(permIDs) > 0 {
			perms, err := svc.authorizationRepo.GetPermissionsByIDs(ctx, permIDs)
			if err == nil {
				role.Permissions = perms
			}
		}
	}

	return role, nil
}
