package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListPermissionsInput struct {
	Keyword string
	Domains []string
}

// ListPermissions retrieves all permission paths, optionally filtered by domain list.
func (r *AuthorizationRepository) ListPermissions(ctx context.Context, input ListPermissionsInput) ([]*models.Permission, error) {
	filter := bson.M{}

	if input.Keyword != "" {
		filter["name"] = bson.M{
			"$regex":   input.Keyword,
			"$options": "i",
		}
	}

	if len(input.Domains) > 0 {
		filter["domain"] = bson.M{"$in": input.Domains}
	}

	opts := options.Find().SetHint(IdxPermissionDomain)

	permissions, err := r.permissionCol.Find(ctx, filter, opts, nil)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}
