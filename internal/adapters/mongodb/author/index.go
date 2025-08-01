package author

import (
	"context"
	"log"

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

// ensureIndexes sets up all necessary indexes for the authorization collections.
func ensureIndexes(repo *AuthorizationRepository) {
	createIndex(repo.PermissionPathCol, bson.D{{Key: "path", Value: 1}}, IdxPermissionPath, true)
	createIndex(repo.PermissionCol, bson.D{
		{Key: "resource", Value: 1},
		{Key: "action", Value: 1},
		{Key: "domain", Value: 1},
	}, IdxPermissionResourceAction, true)
	createIndex(repo.PermissionCol, bson.D{{Key: "domain", Value: 1}}, IdxPermissionDomain, false)

	createIndex(repo.AssignmentCol, bson.D{{Key: "user_id", Value: 1}}, IdxAssignmentUserId, false)
}

// createIndex is a helper to create a MongoDB index.
func createIndex(col *mongo.Collection, keys bson.D, name string, unique bool) {
	_, err := col.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetName(name).SetUnique(unique),
	})
	if err != nil {
		log.Printf("Failed to create index %s: %v", name, err)
	}
}
