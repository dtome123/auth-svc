package auth

import (
	"context"
	"fmt"
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
	if svc.enforcer == nil {
		return false, fmt.Errorf("casbin enforcer is not initialized")
	}

	roles, err := svc.getUserRoles(ctx, input.UserID)
	if err != nil {
		return false, err
	}
	if len(roles) == 0 {
		return false, nil
	}

	return svc.enforcer.EnforceRoles(roles, input.Resource, input.Action)
}
