package author

import (
	"auth-svc/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAssignmentsByRoleID retrieves assignment documents containing a given role ID.
func (repo *impl) GetAssignmentsByRoleID(ctx context.Context, roleID string) ([]*domain.Assignment, error) {
	objID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return nil, nil
	}

	filter := bson.M{
		"role_id": objID,
	}
	docs, err := repo.assignmentCol.Find(ctx, filter, nil, nil)
	if err != nil {
		return nil, err
	}

	assignments := make([]*domain.Assignment, len(docs))
	for i, d := range docs {
		assignments[i] = d.ToDomain()
	}

	return assignments, nil
}
