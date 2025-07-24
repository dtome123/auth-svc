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

type CheckOutput struct {
	Allowed    bool
	StatusCode int32
	Message    string
}

func (svc *AuthorizationService) Check(ctx context.Context, req CheckInput) (CheckOutput, error) {

	claims, err := jwtutils.Extract(req.AccessToken)
	if err != nil {
		return checkInternalServerError(err), err
	}
	userID := claims.Get("sub").AsString()
	deviceID := claims.Get("device_id").AsString()

	// Check cache
	if svc.cfg.Caching.Enable {
		if allowed, found, err := svc.authorCache.GetPermissionCheckResult(ctx, userID, req.FullMethod); err == nil && found && allowed {
			return checkOK(), nil
		}
	}

	permission, err := svc.authorizationRepo.GetPermissionByPath(ctx, req.FullMethod)
	if err != nil {
		return checkInternalServerError(err), err
	}

	// Public route: always allow
	if permission.Type == types.RouteScopePublic {
		svc.cachePermissionCheckResult(ctx, userID, req.FullMethod, true)
		return checkOK(), nil
	}

	if req.AccessToken == "" {
		return checkUnauthorized("token is required"), fmt.Errorf("token is required")
	}

	// Validate session and token
	if err := svc.validateSessionAndToken(ctx, userID, deviceID, req.AccessToken); err != nil {
		return checkUnauthorized("invalid token"), err
	}

	// Check if user has required permission
	hasPermission, err := svc.userHasPermission(ctx, userID, permission)
	if err != nil {
		return checkInternalServerError(err), err
	}

	if !hasPermission {
		return checkForbidden(), nil
	}

	// cache miss
	svc.cachePermissionCheckResult(ctx, userID, req.FullMethod, true)
	return checkOK(), nil
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

	if _, err := svc.serverVerifier.Verify(accessToken); err != nil {
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
