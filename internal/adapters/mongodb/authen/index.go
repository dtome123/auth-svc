package authen

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	IdxSessionUserIdDeviceId = "_idx_session_user_id_device_id_"
	IdxSessionTTL            = "_idx_session_ttl_"
)

func GetSessionIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "device_id", Value: 1},
			},
			Options: options.Index().SetName(IdxSessionUserIdDeviceId).SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "ttl", Value: 1},
			},
			Options: options.Index().SetName(IdxSessionTTL).SetExpireAfterSeconds(0),
		},
	}
}
