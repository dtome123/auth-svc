package author

import (
	"auth-svc/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// ListAllAssignments retrieves all user assignment documents from the database.
func (repo *AuthorizationRepository) ListAllAssignments(ctx context.Context) ([]*domain.Assignment, error) {
	docs, err := repo.assignmentCol.Find(ctx, bson.M{}, nil, nil)
	if err != nil {
		return nil, err
	}

	assignments := make([]*domain.Assignment, len(docs))
	for i, d := range docs {
		assignments[i] = d.ToDomain()
	}

	return assignments, nil
}
