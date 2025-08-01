package auth

import (
	"context"
)

type HasPermissionInput struct {
	UserID   string `json:"user_id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type HasPermissionOutput struct {
	Allowed bool `json:"allowed"`
}

func (svc *AuthorizationService) HasPermission(ctx context.Context, input HasPermissionInput) (bool, error) {
	permissions, err := svc.authorizationRepo.GetPermissionsByUserID(ctx, input.UserID)
	if err != nil {
		return false, err
	}

	for _, permission := range permissions {
		if permission.Resource == input.Resource && permission.Action == input.Action {
			return true, nil
		}
	}

	return false, nil
}
