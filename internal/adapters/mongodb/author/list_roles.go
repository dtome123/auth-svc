package author

import (
	"auth-svc/internal/models"
	"auth-svc/internal/utils"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListRoleInput struct {
	Keyword    string
	Scope      map[string][]string
	PageNumber int64 // starts from 1, if 0 treat as 1
	PageLimit  int64 // max number of items per page
}

// ListRoles retrieves roles from the database with pagination support by page number and page limit.
func (repo *AuthorizationRepository) ListRoles(ctx context.Context, input ListRoleInput) ([]*models.Role, error) {
	skip, limit := utils.PaginationToSkipLimit(input.PageNumber, input.PageLimit)

	findOpts := options.Find()
	if limit > 0 {
		findOpts.SetLimit(limit)
	}
	if skip > 0 {
		findOpts.SetSkip(skip)
	}

	filter := bson.M{}
	if input.Keyword != "" {
		filter["name"] = bson.M{
			"$regex":   input.Keyword,
			"$options": "i", // case-insensitive
		}
	}

	if input.Scope != nil {
		for k, v := range input.Scope {
			if len(v) == 0 {
				continue
			}
			filter[k] = bson.M{"$in": v}
		}
	}

	roles, err := repo.roleCol.Find(ctx, filter, findOpts, nil)
	if err != nil {
		return nil, err
	}

	return roles, nil
}
