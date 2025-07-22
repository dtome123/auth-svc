package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo *AuthorizationRepository) GetRole(ctx context.Context, id string) (models.Role, error) {
	var role models.Role
	err := repo.RoleCol.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
	return role, err
}
