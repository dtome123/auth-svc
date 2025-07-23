package author

import (
	"auth-svc/internal/models"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthorizationCacheRepository struct {
	Redis *redis.Client
	TTL   time.Duration
}

func NewAuthorizationCacheRepository(redisClient *redis.Client, ttl string) *AuthorizationCacheRepository {

	dur, err := time.ParseDuration(ttl)
	if err != nil {
		panic(err)
	}

	return &AuthorizationCacheRepository{
		Redis: redisClient,
		TTL:   dur,
	}
}

//
// === PERMISSIONS CACHE ===
//

func (repo *AuthorizationCacheRepository) GetPermissions(ctx context.Context, userID string) ([]models.Permission, error) {
	var permissions []models.Permission

	cacheKey := repo.buildPermissionListCacheKey(userID)
	cached, err := repo.Redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(cached), &permissions); err != nil {
		return nil, err
	}
	return permissions, nil
}

func (repo *AuthorizationCacheRepository) SetPermissions(ctx context.Context, userID string, permissions []models.Permission) error {
	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	cacheKey := repo.buildPermissionListCacheKey(userID)
	return repo.Redis.Set(ctx, cacheKey, data, repo.TTL).Err()
}

func (repo *AuthorizationCacheRepository) InvalidatePermissions(ctx context.Context, userID string) error {
	cacheKey := repo.buildPermissionListCacheKey(userID)
	return repo.Redis.Del(ctx, cacheKey).Err()
}

func (repo *AuthorizationCacheRepository) ClearUserPermissions(ctx context.Context, userID string) error {
	cacheKey := repo.buildPermissionListCacheKey(userID)
	return repo.Redis.Del(ctx, cacheKey).Err()
}

//
// === PERMISSION CHECK RESULT CACHE (true/false) ===
//

func (repo *AuthorizationCacheRepository) GetPermissionCheckResult(ctx context.Context, userID, fullMethod string) (bool, bool, error) {
	cacheKey := repo.buildPermissionCheckKey(userID, fullMethod)
	cached, err := repo.Redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return false, false, nil // cache miss
	}
	if err != nil {
		return false, false, err
	}

	return cached == "1", true, nil
}

func (repo *AuthorizationCacheRepository) SetPermissionCheckResult(ctx context.Context, userID, fullMethod string, allowed bool) error {
	val := "0"
	if allowed {
		val = "1"
	}

	cacheKey := repo.buildPermissionCheckKey(userID, fullMethod)
	return repo.Redis.Set(ctx, cacheKey, val, repo.TTL).Err()
}

func (repo *AuthorizationCacheRepository) InvalidatePermissionCheckResult(ctx context.Context, userID, fullMethod string) error {
	cacheKey := repo.buildPermissionCheckKey(userID, fullMethod)
	return repo.Redis.Del(ctx, cacheKey).Err()
}

func (repo *AuthorizationCacheRepository) ClearUserPermissionCheck(ctx context.Context, userID string) error {
	cacheKey := repo.buildPermissionCheckKey(userID, "*")
	return repo.Redis.Del(ctx, cacheKey).Err()
}

//
// === KEY BUILDERS ===
//

func (repo *AuthorizationCacheRepository) buildPermissionListCacheKey(userID string) string {
	return "permissions:user:" + userID
}

func (repo *AuthorizationCacheRepository) buildPermissionCheckKey(userID, fullMethod string) string {
	return "permission_result:user:" + userID + ":method:" + fullMethod
}
