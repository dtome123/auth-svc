package author

import (
	mongoModels "auth-svc/internal/adapters/mongodb/models"
	"auth-svc/internal/domain"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BatchUpsertPermissions upserts a list of permissions into the database.
func (repo *AuthorizationRepository) BatchUpsertPermissions(ctx context.Context, permissions []*domain.Permission) error {
	if len(permissions) == 0 {
		return nil
	}

	writeModels := make([]mongo.WriteModel, 0, len(permissions))
	for _, perm := range permissions {
		doc := mongoModels.PermissionFromDomain(perm)
		filter := bson.M{
			"_id": doc.ID,
		}

		writeModels = append(writeModels, mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(bson.M{"$set": doc}).
			SetUpsert(true),
		)
	}

	_, err := repo.permissionCol.BulkWrite(ctx, writeModels, &options.BulkWriteOptions{})
	if err != nil {
		log.Printf("Failed to bulk upsert permissions: %v", err)
	}
	return err
}
