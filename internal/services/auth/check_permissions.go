package auth

import (
	"auth-svc/internal/models"
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type CheckPermissionsInput struct {
	UserID          string                  `json:"user_id"`
	ActionResources []models.ActionResource `json:"action_resources"`
}

type CheckPermissionsOutput struct {
	Allowed bool `json:"allowed"`
}

func (svc *AuthorizationService) CheckPermissions(ctx context.Context, input CheckPermissionsInput) (map[string]bool, error) {
	permissions, err := svc.authorizationRepo.GetPermissionsByUserID(ctx, input.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get permissions")
	}

	// Build a set of requested keys like "resource:action"
	requested := make(map[string]bool, len(input.ActionResources))
	for _, ar := range input.ActionResources {
		key := fmt.Sprintf("%s:%s", ar.Resource, ar.Action)
		requested[key] = true
	}

	result := make(map[string]bool, len(requested))
	// Mark allowed keys as true
	for _, p := range permissions {
		key := fmt.Sprintf("%s:%s", p.Resource, p.Action)
		if requested[key] {
			result[key] = true
		}
	}

	// For keys not in result, default to false
	for key := range requested {
		if _, found := result[key]; !found {
			result[key] = false
		}
	}

	return result, nil
}
