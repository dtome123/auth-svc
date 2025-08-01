package authen

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	IdxSessionUserIdDeviceId = "_idx_session_user_id_device_id_"
	IdxSessionTTL            = "_idx_session_ttl_"
)

// indexingSessionCol creates all necessary indexes on the session collection.
func indexingSessionCol(ctx context.Context, col *mongo.Collection) {
	createIndex(ctx, col, bson.D{
		{Key: "user_id", Value: 1},
		{Key: "device_id", Value: 1},
	}, IdxSessionUserIdDeviceId, true)

	createIndex(ctx, col, bson.D{
		{Key: "ttl", Value: 1},
	}, IdxSessionTTL, false, options.Index().SetExpireAfterSeconds(0))
}

// createIndex is a helper to create a MongoDB index with optional index options.
func createIndex(ctx context.Context, col *mongo.Collection, keys bson.D, name string, unique bool, opts ...*options.IndexOptions) {
	indexOpts := options.Index().SetName(name).SetUnique(unique)
	if len(opts) > 0 && opts[0] != nil {
		// Merge options if extra options passed
		if opts[0].ExpireAfterSeconds != nil {
			indexOpts.SetExpireAfterSeconds(*opts[0].ExpireAfterSeconds)
		}
		// Add here more options merging if needed
	}

	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    keys,
		Options: indexOpts,
	})
	if err != nil {
		log.Printf("Failed to create index %s: %v", name, err)
	}
}
