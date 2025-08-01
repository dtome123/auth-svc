package author

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthorizationRepository struct {
	AssignmentCol     *mongo.Collection
	RoleCol           *mongo.Collection
	PermissionCol     *mongo.Collection
	PermissionPathCol *mongo.Collection
}

// NewAuthorizationRepository initializes collections and indexes for authorization.
func NewAuthorizationRepository(db *mongo.Database) *AuthorizationRepository {
	repo := &AuthorizationRepository{
		AssignmentCol:     db.Collection("assignments"),
		RoleCol:           db.Collection("roles"),
		PermissionCol:     db.Collection("permissions"),
		PermissionPathCol: db.Collection("path_permissions"),
	}

	// Setup indexes for faster access
	ensureIndexes(repo)

	return repo
}
