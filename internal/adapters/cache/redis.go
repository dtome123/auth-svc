package cache

import (
	"auth-svc/internal/domain"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisAuthorizationCacheRepository struct {
	Redis       *redis.Client
	durationTTL time.Duration
}

func NewRedisAuthorizationCacheRepository(redisClient *redis.Client, durationTTL string) AuthorizationCacheRepository {
	dur, err := time.ParseDuration(durationTTL)
	if err != nil {
		panic(err)
	}

	return &RedisAuthorizationCacheRepository{
		Redis:       redisClient,
		durationTTL: dur,
	}
}

//
// === PERMISSIONS CACHE ===
//

func (repo *RedisAuthorizationCacheRepository) GetPermissions(ctx context.Context, userID string) ([]domain.ResourceAction, error) {
	var permissions []domain.ResourceAction

	cacheKey := BuildPermissionListCacheKey(userID)
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

func (repo *RedisAuthorizationCacheRepository) SetPermissions(ctx context.Context, userID string, permissions []domain.ResourceAction) error {
	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	cacheKey := BuildPermissionListCacheKey(userID)
	return repo.Redis.Set(ctx, cacheKey, data, repo.durationTTL).Err()
}

func (repo *RedisAuthorizationCacheRepository) InvalidatePermissions(ctx context.Context, userID string) error {
	cacheKey := BuildPermissionListCacheKey(userID)
	return repo.Redis.Del(ctx, cacheKey).Err()
}

func (repo *RedisAuthorizationCacheRepository) ClearUserPermissions(ctx context.Context, userID string) error {
	cacheKey := BuildPermissionListCacheKey(userID)
	return repo.Redis.Del(ctx, cacheKey).Err()
}

//
// === PERMISSION CHECK RESULT CACHE (true/false) ===
//

func (repo *RedisAuthorizationCacheRepository) GetPermissionCheckResult(ctx context.Context, userID, fullMethod string) (bool, bool, error) {
	cacheKey := BuildPermissionCheckKey(userID, fullMethod)
	cached, err := repo.Redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return false, false, nil // cache miss
	}
	if err != nil {
		return false, false, err
	}

	return cached == "1", true, nil
}

func (repo *RedisAuthorizationCacheRepository) SetPermissionCheckResult(ctx context.Context, userID, fullMethod string, allowed bool) error {
	val := "0"
	if allowed {
		val = "1"
	}

	cacheKey := BuildPermissionCheckKey(userID, fullMethod)
	return repo.Redis.Set(ctx, cacheKey, val, repo.durationTTL).Err()
}

func (repo *RedisAuthorizationCacheRepository) InvalidatePermissionCheckResult(ctx context.Context, userID, fullMethod string) error {
	cacheKey := BuildPermissionCheckKey(userID, fullMethod)
	return repo.Redis.Del(ctx, cacheKey).Err()
}

func (repo *RedisAuthorizationCacheRepository) ClearUserPermissionCheck(ctx context.Context, userID string) error {
	pattern := BuildPermissionCheckKey(userID, "*")
	var cursor uint64
	for {
		keys, nextCursor, err := repo.Redis.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := repo.Redis.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

func (repo *RedisAuthorizationCacheRepository) ClearAllPermissionsAndChecks(ctx context.Context) error {
	patterns := []string{KeyPermissionsPrefix + "*", KeyPermissionCheckPrefix + "*"}
	for _, pattern := range patterns {
		var cursor uint64
		for {
			keys, nextCursor, err := repo.Redis.Scan(ctx, cursor, pattern, 100).Result()
			if err != nil {
				return err
			}
			if len(keys) > 0 {
				if err := repo.Redis.Del(ctx, keys...).Err(); err != nil {
					return err
				}
			}
			cursor = nextCursor
			if cursor == 0 {
				break
			}
		}
	}
	return nil
}
