package author

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	IdxPermissionDomainResource = "_idx_permission_domain_resource_"
	IdxAssignmentUserId         = "_idx_permission_user_id_"
)

func GetPermissionIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "domain", Value: 1},
				{Key: "resource", Value: 1},
			},
			Options: options.Index().SetName(IdxPermissionDomainResource).SetUnique(true),
		},
	}
}

func GetAssignmentIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetName(IdxAssignmentUserId).SetUnique(true),
		},
	}
}

func GetRoleIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
			},
			Options: options.Index().SetName("_idx_role_name_").SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "code", Value: 1},
			},
			Options: options.Index().SetName("_idx_role_code_").SetUnique(true).SetSparse(true),
		},
	}
}
