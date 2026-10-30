package cache

import (
	"auth-svc/internal/domain"
	"context"
)

type AuthorizationCacheRepository interface {
	GetPermissions(ctx context.Context, userID string) ([]domain.ResourceAction, error)
	SetPermissions(ctx context.Context, userID string, permissions []domain.ResourceAction) error
	InvalidatePermissions(ctx context.Context, userID string) error
	ClearUserPermissions(ctx context.Context, userID string) error
	GetPermissionCheckResult(ctx context.Context, userID, fullMethod string) (allowed bool, found bool, err error)
	SetPermissionCheckResult(ctx context.Context, userID, fullMethod string, allowed bool) error
	InvalidatePermissionCheckResult(ctx context.Context, userID, fullMethod string) error
	ClearUserPermissionCheck(ctx context.Context, userID string) error
	ClearAllPermissionsAndChecks(ctx context.Context) error
}
