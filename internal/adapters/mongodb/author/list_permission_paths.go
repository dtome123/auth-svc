package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListPermissionPathsInput struct {
	Keyword string
	Domains []string
}

// ListPermissionPaths retrieves all permission paths, optionally filtered by domain list.
func (r *AuthorizationRepository) ListPermissionPaths(ctx context.Context, input ListPermissionPathsInput) ([]*models.PermissionPath, error) {
	filter := bson.M{}

	if input.Keyword != "" {
		filter["path"] = bson.M{
			"$regex":   input.Keyword,
			"$options": "i",
		}
	}

	if len(input.Domains) > 0 {
		filter["domain"] = bson.M{"$in": input.Domains}
	}

	opts := options.Find().SetHint(IdxPermissionPath)

	permissions, err := r.permissionPathCol.Find(ctx, filter, opts, nil)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}
