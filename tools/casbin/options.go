package casbin

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// Option configures an IAM enforcer instance during construction.
type Option func(*impl)

// WithRedis configures Redis Pub/Sub for cross-instance policy synchronization.
func WithRedis(client *redis.Client) Option {
	return func(i *impl) {
		if client != nil {
			i.sync = newSyncer(client, i.cfg.SyncChannel)
		}
	}
}

// WithSyncChannel sets a custom Redis Pub/Sub channel for policy synchronization.
func WithSyncChannel(channel string) Option {
	return func(i *impl) {
		if channel != "" {
			i.cfg.SyncChannel = channel
		}
	}
}

// WithSyncTimeout sets the maximum timeout duration for reload operations triggered by sync events.
func WithSyncTimeout(timeout time.Duration) Option {
	return func(i *impl) {
		if timeout > 0 {
			i.cfg.SyncTimeout = timeout
		}
	}
}
