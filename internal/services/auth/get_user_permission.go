package auth

import (
	"auth-svc/internal/domain"
	"context"
)

func (svc *AuthorizationService) GetUserPermissions(ctx context.Context, userID string) ([]domain.ResourceAction, error) {

	if svc.cfg.Caching.Enable {
		perms, err := svc.authorCache.GetPermissions(ctx, userID)
		if err != nil {
			return nil, err
		}
		if perms != nil {
			return perms, nil
		}
	}

	actionResources, err := svc.authorizationRepo.GetPermissionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if svc.cfg.Caching.Enable {
		if err := svc.authorCache.SetPermissions(ctx, userID, actionResources); err != nil {
			return nil, err
		}
	}

	return actionResources, nil
}
