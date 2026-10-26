package cache

const (
	KeyPermissionsPrefix     = "permissions:user:"
	KeyPermissionCheckPrefix = "permission_result:user:"
)

func BuildPermissionListCacheKey(userID string) string {
	return KeyPermissionsPrefix + userID
}

func BuildPermissionCheckKey(userID, fullMethod string) string {
	return KeyPermissionCheckPrefix + userID + ":method:" + fullMethod
}

func BuildPermissionCheckUserPrefix(userID string) string {
	return KeyPermissionCheckPrefix + userID + ":"
}
