package casbin

import "time"

// Config holds the configuration for constructing a new Enforcer.
// All fields have sensible defaults and are optional.
type Config struct {
	// SyncChannel is the Redis Pub/Sub channel name used to broadcast
	// policy-reload events across service instances.
	// Default: "casbin:policy:sync"
	SyncChannel string `yaml:"sync_channel" json:"sync_channel"`

	// SyncTimeout is the maximum duration allowed for a reload triggered
	// by a remote sync event.
	// Default: 10s
	SyncTimeout time.Duration `yaml:"sync_timeout" json:"sync_timeout"`
}

// DefaultConfig returns a Config populated with production-ready defaults.
func DefaultConfig() Config {
	return Config{
		SyncChannel: "casbin:policy:sync",
		SyncTimeout: 10 * time.Second,
	}
}

func (c Config) withDefaults() Config {
	if c.SyncChannel == "" {
		c.SyncChannel = DefaultConfig().SyncChannel
	}
	if c.SyncTimeout <= 0 {
		c.SyncTimeout = DefaultConfig().SyncTimeout
	}
	return c
}
