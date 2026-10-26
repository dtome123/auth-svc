package cache

import (
	"auth-svc/internal/domain"
	"context"
)

type NoOpAuthorizationCacheRepository struct{}

func NewNoOpAuthorizationCacheRepository() AuthorizationCacheRepository {
	return &NoOpAuthorizationCacheRepository{}
}

func (n *NoOpAuthorizationCacheRepository) GetPermissions(ctx context.Context, userID string) ([]*domain.Permission, error) {
	return nil, nil
}

func (n *NoOpAuthorizationCacheRepository) SetPermissions(ctx context.Context, userID string, permissions []*domain.Permission) error {
	return nil
}

func (n *NoOpAuthorizationCacheRepository) InvalidatePermissions(ctx context.Context, userID string) error {
	return nil
}

func (n *NoOpAuthorizationCacheRepository) ClearUserPermissions(ctx context.Context, userID string) error {
	return nil
}

func (n *NoOpAuthorizationCacheRepository) GetPermissionCheckResult(ctx context.Context, userID, fullMethod string) (bool, bool, error) {
	return false, false, nil
}

func (n *NoOpAuthorizationCacheRepository) SetPermissionCheckResult(ctx context.Context, userID, fullMethod string, allowed bool) error {
	return nil
}

func (n *NoOpAuthorizationCacheRepository) InvalidatePermissionCheckResult(ctx context.Context, userID, fullMethod string) error {
	return nil
}

func (n *NoOpAuthorizationCacheRepository) ClearUserPermissionCheck(ctx context.Context, userID string) error {
	return nil
}

func (n *NoOpAuthorizationCacheRepository) ClearAllPermissionsAndChecks(ctx context.Context) error {
	return nil
}
