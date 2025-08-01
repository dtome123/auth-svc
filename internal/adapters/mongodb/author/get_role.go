package author

import (
	"auth-svc/internal/models"
	"context"

	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetRole retrieves a role by its ID (hex string).
func (repo *AuthorizationRepository) GetRole(ctx context.Context, id string) (*models.Role, error) {
	var role *models.Role

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid role ID format")
	}

	err = repo.RoleCol.FindOne(ctx, bson.M{"_id": objID}).Decode(&role)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}
	return role, nil
}
