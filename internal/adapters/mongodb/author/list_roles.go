package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type ListRoleInput struct {
}

func (repo *AuthorizationRepository) ListRoles(ctx context.Context, input ListRoleInput) ([]models.Role, error) {
	var roles []models.Role
	cursor, err := repo.RoleCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}
