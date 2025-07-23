package auth

import (
	"auth-svc/internal/models"
	"auth-svc/internal/types"
	"auth-svc/internal/utils"
	"context"
	"fmt"

	"github.com/dtome123/auth-sdk/jwtutils"
)

type CheckInput struct {
	AccessToken string
	FullMethod  string
}

func (svc *AuthorizationService) Check(ctx context.Context, req CheckInput) (bool, error) {
	if req.AccessToken == "" {
		return false, fmt.Errorf("token is required")
	}

	claims, err := jwtutils.Extract(req.AccessToken)
	if err != nil {
		return false, err
	}
	userID := claims.Get("sub").AsString()
	deviceID := claims.Get("device_id").AsString()

	// Check cache
	if svc.cfg.Caching.Enable {
		if allowed, found, err := svc.authorCache.GetPermissionCheckResult(ctx, userID, req.FullMethod); err == nil && found && allowed {
			return true, nil
		}
	}

	permission, err := svc.authorizationRepo.GetPermissionByPath(ctx, req.FullMethod)
	if err != nil {
		return false, err
	}

	// Public route: always allow
	if permission.Type == types.RouteScopePublic {
		svc.cachePermissionCheckResult(ctx, userID, req.FullMethod, true)
		return true, nil
	}

	// Validate session and token
	if err := svc.validateSessionAndToken(ctx, userID, deviceID, req.AccessToken); err != nil {
		return false, err
	}

	// Check if user has required permission
	hasPermission, err := svc.userHasPermission(ctx, userID, permission)
	if err != nil {
		return false, err
	}

	if hasPermission {
		svc.cachePermissionCheckResult(ctx, userID, req.FullMethod, true)
	}

	return hasPermission, nil
}

func (svc *AuthorizationService) validateSessionAndToken(ctx context.Context, userID, deviceID, accessToken string) error {
	accessTokenHash := utils.HashSHA256(accessToken)
	session, err := svc.authenticationRepo.GetSession(ctx, userID, deviceID)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("no session found for this token")
	}
	if session.AccessTokenHash != accessTokenHash {
		return fmt.Errorf("invalid token")
	}

	if _, err := jwtutils.VerifyJWTWithRS256(accessToken, svc.cfg.Cert.PublicKeyPath); err != nil {
		return err
	}

	return nil
}

func (svc *AuthorizationService) userHasPermission(ctx context.Context, userID string, permission models.PermissionPath) (bool, error) {
	var perms []models.Permission
	var err error

	// Try to get from cache
	if svc.cfg.Caching.Enable {
		if cachedPerms, err := svc.authorCache.GetPermissions(ctx, userID); err == nil && len(cachedPerms) > 0 {
			perms = cachedPerms
		}
	}

	// If not found in cache, load from DB
	if perms == nil {
		perms, err = svc.GetUserPermissions(ctx, userID)
		if err != nil {
			return false, err
		}

		// Save to cache
		if svc.cfg.Caching.Enable {
			if err := svc.authorCache.SetPermissions(ctx, userID, perms); err != nil {
				fmt.Printf("failed to cache user permissions: %v\n", err)
			}
		}
	}

	requireAction := models.ActionResource{
		Resource: permission.Resource,
		Action:   permission.Action,
	}

	// Check if user has required permission
	for _, p := range perms {
		if p.Resource == requireAction.Resource && p.Action == requireAction.Action {
			return true, nil
		}
		for _, implied := range p.ImpliedByActions {
			if implied.Resource == requireAction.Resource && implied.Action == requireAction.Action {
				return true, nil
			}
		}
	}

	return false, nil
}

func (svc *AuthorizationService) cachePermissionCheckResult(ctx context.Context, userID, fullMethod string, result bool) {
	if !svc.cfg.Caching.Enable {
		return
	}
	if err := svc.authorCache.SetPermissionCheckResult(ctx, userID, fullMethod, result); err != nil {
		fmt.Printf("failed to cache permission check result: %v \n", err)
	}
}
