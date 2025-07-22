package author

import (
	"auth-svc/internal/models"
	"context"
)

func (svc *AuthorizationService) GetUserPermissions(ctx context.Context, userID string) ([]models.Permission, error) {
	perms, err := svc.Cache.GetPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}
	if perms != nil {
		return perms, nil
	}

	perms, err = svc.Repo.GetPermissionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := svc.Cache.SetPermissions(ctx, userID, perms); err != nil {
		return nil, err
	}
	return perms, nil
}
