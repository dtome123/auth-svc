package auth

import (
	"auth-svc/internal/models"
	"context"
)

func (svc *AuthorizationService) GetUserPermissions(ctx context.Context, userID string) ([]models.Permission, error) {

	if svc.cfg.Caching.Enable {
		perms, err := svc.authorCache.GetPermissions(ctx, userID)
		if err != nil {
			return nil, err
		}
		if perms != nil {
			return perms, nil
		}
	}

	perms, err := svc.authorizationRepo.GetPermissionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if svc.cfg.Caching.Enable {
		if err := svc.authorCache.SetPermissions(ctx, userID, perms); err != nil {
			return nil, err
		}
	}

	return perms, nil
}
