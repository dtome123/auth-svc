package author

import (
	"auth-svc/internal/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// BatchUpsertPermissions upserts a list of permissions into the database.
func (repo *AuthorizationRepository) BatchUpsertPermissions(ctx context.Context, permissions []*models.Permission) error {
	if len(permissions) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, 0, len(permissions))
	for _, perm := range permissions {
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{
				"domain":   perm.Domain,
				"resource": perm.Resource,
				"action":   perm.Action,
			}).
			SetUpdate(bson.M{"$set": perm}).
			SetUpsert(true),
		)
	}

	_, err := repo.PermissionCol.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("Failed to bulk upsert permissions: %v", err)
	}
	return err
}
