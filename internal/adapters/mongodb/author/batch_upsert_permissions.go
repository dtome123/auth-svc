package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (repo *AuthorizationRepository) BatchUpsertPermissions(ctx context.Context, permissions []*models.Permission) error {

	var models []mongo.WriteModel
	for _, permission := range permissions {
		models = append(models, mongo.NewUpdateOneModel().SetFilter(
			bson.M{
				"domain":   permission.Domain,
				"resource": permission.Resource,
				"action":   permission.Action,
			}).
			SetUpsert(true).
			SetUpdate(bson.M{"$set": permission}),
		)
	}
	_, err := repo.PermissionCol.BulkWrite(ctx, models)
	return err
}
