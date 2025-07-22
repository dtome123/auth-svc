package author

import (
	"auth-svc/internal/models"
	"context"
	"fmt"
	"strings"
)

// UserHasPermission checks if the user has permission to call fullMethod
func (svc *AuthorizationService) UserHasPermission(ctx context.Context, userID, fullMethod string) (bool, error) {
	// Check cache first; if cached result is true, return immediately
	if cachedResult, found, err := svc.Cache.GetPermissionCheckResult(ctx, userID, fullMethod); err == nil && found && cachedResult {
		return true, nil
	}

	// Load user's permissions from DB or cache
	perms, err := svc.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	permission, err := svc.Repo.GetPermissionByPath(ctx, fullMethod)
	if err != nil {
		return false, err
	}

	requireAction := models.ActionResource{
		Resource: permission.Resource,
		Action:   permission.Action,
	}

	hasPermission := false
	for _, p := range perms {
		if p.Resource == requireAction.Resource && p.Action == requireAction.Action {
			hasPermission = true
			return true, nil
		}

		for _, i := range p.ImpliedByActions {
			if i.Resource == requireAction.Resource && i.Action == requireAction.Action {
				hasPermission = true
				return true, nil
			}
		}
	}

	if hasPermission {
		if err := svc.Cache.SetPermissionCheckResult(ctx, userID, fullMethod, true); err != nil {
			fmt.Printf("failed to cache permission check result: %v \n", err)
		}
	}

	return hasPermission, nil

}

// pathMatches supports wildcard "*" at the end of permission path for prefix matching
func pathMatches(permPath, methodPath string) bool {
	if strings.HasSuffix(permPath, "/*") {
		prefix := strings.TrimSuffix(permPath, "/*")
		return strings.HasPrefix(methodPath, prefix)
	}
	return permPath == methodPath
}
