package author

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	IdxPermissionPath = "_idx_permission_path_"

	IdxPermissionResourceAction = "_idx_permission_resource_action_"

	IdxAssignmentUserId = "_idx_permission_user_id_"
)

func indexingPermissionPathCol(col *mongo.Collection) {
	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"path": 1,
		},
		Options: options.Index().SetName(IdxPermissionPath).SetUnique(true),
	})
}

func indexingPermissionCol(col *mongo.Collection) {
	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"resource": 1,
			"action":   1,
			"domain":   1,
		},
		Options: options.Index().SetName(IdxPermissionResourceAction).SetUnique(true),
	})
}

func indexingAssignmentCol(col *mongo.Collection) {
	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"user_id": 1,
		},
		Options: options.Index().SetName(IdxAssignmentUserId),
	})
}
