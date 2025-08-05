package models

import (
	"auth-svc/internal/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Permission struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Domain           string             `bson:"domain" json:"domain"`
	Description      string             `bson:"description" json:"description"`
	Resource         string             `bson:"resource" json:"resource"`
	Action           string             `bson:"action" json:"action"`
	ImpliedByActions []ActionResource   `bson:"implied_actions,omitempty" json:"implied_actions,omitempty"`
}

func (Permission) CollectionName() string {
	return "permissions"
}

type ActionResource struct {
	Resource string `bson:"resource" json:"resource"`
	Action   string `bson:"action" json:"action"`
}

type PermissionPath struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Domain   string             `bson:"domain" json:"domain"`
	Path     string             `bson:"path" json:"path"`
	Resource string             `bson:"resource" json:"resource"`
	Action   string             `bson:"action" json:"action"`
	Type     types.RouteScope   `bson:"type" json:"type"`
}

func (PermissionPath) CollectionName() string {
	return "permission_paths"
}