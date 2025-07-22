package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo *AuthorizationRepository) UpdateRole(ctx context.Context, role *models.Role) error {
	_, err := repo.RoleCol.UpdateOne(ctx, bson.M{"_id": role.ID}, bson.M{"$set": role})
	return err
}
