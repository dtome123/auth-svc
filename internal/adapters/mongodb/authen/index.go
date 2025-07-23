package authen

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	IdxSessionUserIdDeviceId = "_idx_session_user_id_device_id_"
	IdxSessionTTL            = "_idx_session_ttl_"
)

func indexingSessionCol(col *mongo.Collection) {

	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"user_id":   1,
			"device_id": 1,
		},
		Options: options.Index().SetName(IdxSessionUserIdDeviceId).SetUnique(true),
	})

	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"ttl": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(0).SetName(IdxSessionTTL),
	})
}
