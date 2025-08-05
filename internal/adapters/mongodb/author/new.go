package author

import (
	"auth-svc/internal/models"

	mongodb "github.com/dtome123/go-mongo-generic"
)

type AuthorizationRepository struct {
	assignmentCol     mongodb.Collection[models.Assignment]
	roleCol           mongodb.Collection[models.Role]
	permissionCol     mongodb.Collection[models.Permission]
	permissionPathCol mongodb.Collection[models.PermissionPath]
}

// NewAuthorizationRepository initializes collections and indexes for authorization.
func NewAuthorizationRepository(db *mongodb.Database) *AuthorizationRepository {

	repo := &AuthorizationRepository{
		assignmentCol:     mongodb.NewCollection[models.Assignment](db),
		roleCol:           mongodb.NewCollection[models.Role](db),
		permissionCol:     mongodb.NewCollection[models.Permission](db),
		permissionPathCol: mongodb.NewCollection[models.PermissionPath](db),
	}

	repo.assignmentCol.EnsureIndexes(GetAssignmentIndexes())
	repo.permissionCol.EnsureIndexes(GetPermissionIndexes())
	repo.permissionPathCol.EnsureIndexes(GetPermissionPathIndexes())

	return repo
}
