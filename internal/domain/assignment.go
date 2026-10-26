package domain

// Assignment represents the association of roles to a specific user.
type Assignment struct {
	UserID  string   `json:"user_id"`
	RoleIDs []string `json:"role_id"`
	Roles   []Role   `json:"role,omitempty"`
}
