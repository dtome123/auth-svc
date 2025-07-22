package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Permission struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Domain           string             `bson:"domain" json:"domain"`
	Resource         string             `bson:"resource" json:"resource"`
	Action           string             `bson:"action" json:"action"`
	ImpliedByActions []ActionResource   `bson:"implied_actions,omitempty" json:"implied_actions,omitempty"`
}

type ActionResource struct {
	Resource string `bson:"resource" json:"resource"`
	Action   string `bson:"action" json:"action"`
}

type PermissionPath struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Path     string             `bson:"path" json:"path"`
	Resource string             `bson:"resource" json:"resource"`
	Action   string             `bson:"action" json:"action"`
}
