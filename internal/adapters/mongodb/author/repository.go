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

func NewAuthorizationCacheRepository(db *mongo.Database) *AuthorizationRepository {
	assignmentCol := db.Collection("assignments")
	roleCol := db.Collection("roles")
	permissionCol := db.Collection("permissions")
	pathPermissionCol := db.Collection("path_permissions")

	indexingAssignmentCol(assignmentCol)
	indexingPermissionCol(permissionCol)
	indexingPermissionPathCol(pathPermissionCol)

	return &AuthorizationRepository{
		AssignmentCol:     assignmentCol,
		RoleCol:           roleCol,
		PermissionCol:     permissionCol,
		PermissionPathCol: pathPermissionCol,
	}

}
