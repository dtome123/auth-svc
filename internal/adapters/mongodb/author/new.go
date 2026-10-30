package author

import (
	mongoModels "auth-svc/internal/adapters/mongodb/models"

	mongodb "github.com/dtome123/go-mongo-generic"
)

type impl struct {
	assignmentCol mongodb.Collection[mongoModels.Assignment]
	roleCol       mongodb.Collection[mongoModels.Role]
	permissionCol mongodb.Collection[mongoModels.Permission]
}

// NewAuthorizationRepository initializes collections and indexes for authorization.
func NewAuthorizationRepository(db *mongodb.Database) AuthorizationRepository {

	repo := &impl{
		assignmentCol: mongodb.NewCollection[mongoModels.Assignment](db),
		roleCol:       mongodb.NewCollection[mongoModels.Role](db),
		permissionCol: mongodb.NewCollection[mongoModels.Permission](db),
	}

	repo.assignmentCol.EnsureIndexes(GetAssignmentIndexes())
	repo.roleCol.EnsureIndexes(GetRoleIndexes())
	repo.permissionCol.EnsureIndexes(GetPermissionIndexes())

	return repo
}
