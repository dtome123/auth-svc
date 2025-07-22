package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (repo *AuthorizationRepository) BatchUpsertPermissionPaths(ctx context.Context, paths []*models.PermissionPath) error {

	var models []mongo.WriteModel
	for _, path := range paths {
		models = append(models, mongo.NewUpdateOneModel().SetFilter(
			bson.M{
				"path": path.Path,
			}).
			SetUpsert(true).
			SetUpdate(bson.M{"$set": path}),
		)
	}
	_, err := repo.PermissionPathCol.BulkWrite(ctx, models)
	return err
}
