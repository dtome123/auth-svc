package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name          string               `bson:"name" json:"name"`
	PermissionIDs []primitive.ObjectID `bson:"permission_ids" json:"permission_ids"`
	Permissions   []Permission         `bson:"-" json:"permissions"`
}
