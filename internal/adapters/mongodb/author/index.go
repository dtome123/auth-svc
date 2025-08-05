package author

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	IdxPermissionPath           = "_idx_permission_path_"
	IdxPermissionResourceAction = "_idx_permission_resource_action_"
	IdxPermissionDomain         = "_idx_permission_domain_"
	IdxAssignmentUserId         = "_idx_permission_user_id_"
)

func GetPermissionPathIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "path", Value: 1},
			},
			Options: options.Index().SetName(IdxPermissionPath).SetUnique(true),
		},
	}
}

func GetPermissionIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "resource", Value: 1},
				{Key: "action", Value: 1},
				{Key: "domain", Value: 1},
			},
			Options: options.Index().SetName(IdxPermissionResourceAction).SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "domain", Value: 1},
			},
			Options: options.Index().SetName(IdxPermissionDomain),
		},
	}
}

func GetAssignmentIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetName(IdxAssignmentUserId),
		},
	}
}
