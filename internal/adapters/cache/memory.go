package cache

import (
	"auth-svc/internal/domain"
	"context"
	"strings"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type MemoryAuthorizationCacheRepository struct {
	client *gocache.Cache
}

func NewMemoryAuthorizationCacheRepository(durationTTL string) AuthorizationCacheRepository {
	dur, err := time.ParseDuration(durationTTL)
	if err != nil {
		dur = 5 * time.Minute
	}

	cleanupInterval := dur / 2
	if cleanupInterval < time.Minute {
		cleanupInterval = time.Minute
	}

	return &MemoryAuthorizationCacheRepository{
		client: gocache.New(dur, cleanupInterval),
	}
}

// === PERMISSIONS CACHE ===

func (m *MemoryAuthorizationCacheRepository) GetPermissions(ctx context.Context, userID string) ([]domain.ResourceAction, error) {
	cacheKey := BuildPermissionListCacheKey(userID)
	val, found := m.client.Get(cacheKey)
	if !found {
		return nil, nil
	}

	if permissions, ok := val.([]domain.ResourceAction); ok {
		return permissions, nil
	}

	return nil, nil
}

func (m *MemoryAuthorizationCacheRepository) SetPermissions(ctx context.Context, userID string, permissions []domain.ResourceAction) error {
	cacheKey := BuildPermissionListCacheKey(userID)
	m.client.Set(cacheKey, permissions, gocache.DefaultExpiration)
	return nil
}

func (m *MemoryAuthorizationCacheRepository) InvalidatePermissions(ctx context.Context, userID string) error {
	cacheKey := BuildPermissionListCacheKey(userID)
	m.client.Delete(cacheKey)
	return nil
}

func (m *MemoryAuthorizationCacheRepository) ClearUserPermissions(ctx context.Context, userID string) error {
	cacheKey := BuildPermissionListCacheKey(userID)
	m.client.Delete(cacheKey)
	return nil
}

// === PERMISSION CHECK RESULT CACHE ===

func (m *MemoryAuthorizationCacheRepository) GetPermissionCheckResult(ctx context.Context, userID, fullMethod string) (bool, bool, error) {
	cacheKey := BuildPermissionCheckKey(userID, fullMethod)
	val, found := m.client.Get(cacheKey)
	if !found {
		return false, false, nil
	}

	if allowed, ok := val.(bool); ok {
		return allowed, true, nil
	}

	return false, false, nil
}

func (m *MemoryAuthorizationCacheRepository) SetPermissionCheckResult(ctx context.Context, userID, fullMethod string, allowed bool) error {
	cacheKey := BuildPermissionCheckKey(userID, fullMethod)
	m.client.Set(cacheKey, allowed, gocache.DefaultExpiration)
	return nil
}

func (m *MemoryAuthorizationCacheRepository) InvalidatePermissionCheckResult(ctx context.Context, userID, fullMethod string) error {
	cacheKey := BuildPermissionCheckKey(userID, fullMethod)
	m.client.Delete(cacheKey)
	return nil
}

func (m *MemoryAuthorizationCacheRepository) ClearUserPermissionCheck(ctx context.Context, userID string) error {
	prefix := BuildPermissionCheckUserPrefix(userID)
	for k := range m.client.Items() {
		if strings.HasPrefix(k, prefix) {
			m.client.Delete(k)
		}
	}
	return nil
}

func (m *MemoryAuthorizationCacheRepository) ClearAllPermissionsAndChecks(ctx context.Context) error {
	m.client.Flush()
	return nil
}
