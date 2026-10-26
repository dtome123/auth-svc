package casbin

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// syncer encapsulates the Redis Pub/Sub synchronization logic.
type syncer struct {
	client  *redis.Client
	channel string
	onEvent func(ctx context.Context) error
}

// newSyncer creates a syncer instance. Returns nil if client is nil.
func newSyncer(client *redis.Client, channel string) *syncer {
	if client == nil {
		return nil
	}
	return &syncer{
		client:  client,
		channel: channel,
	}
}

// publish sends a reload notification to the Pub/Sub channel.
func (s *syncer) publish(ctx context.Context) {
	if s == nil {
		return
	}
	if err := s.client.Publish(ctx, s.channel, "reload").Err(); err != nil {
		log.Printf("⚠️ casbin: failed to publish sync event on %q: %v", s.channel, err)
	}
}

// subscribe blocks and listens for reload events. Typically run in a goroutine.
func (s *syncer) subscribe(timeout context.Context, cfg Config) {
	if s == nil || s.onEvent == nil {
		return
	}
	pubsub := s.client.Subscribe(context.Background(), s.channel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for range ch {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.SyncTimeout)
		if err := s.onEvent(ctx); err != nil {
			log.Printf("⚠️ casbin: reload via sync channel %q failed: %v", s.channel, err)
		}
		cancel()
	}
}
