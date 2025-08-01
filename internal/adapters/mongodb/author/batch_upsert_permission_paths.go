package author

import (
	"auth-svc/internal/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// BatchUpsertPermissionPaths upserts a list of permission paths into the database.
func (repo *AuthorizationRepository) BatchUpsertPermissionPaths(ctx context.Context, paths []*models.PermissionPath) error {
	if len(paths) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, 0, len(paths))
	for _, path := range paths {
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"path": path.Path}).
			SetUpdate(bson.M{"$set": path}).
			SetUpsert(true),
		)
	}

	_, err := repo.PermissionPathCol.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("Failed to bulk upsert permission paths: %v", err)
	}
	return err
}
