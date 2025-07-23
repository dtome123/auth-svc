package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Assignment struct {
	UserID  primitive.ObjectID   `bson:"user_id" json:"user_id"`
	RoleIDs []primitive.ObjectID `bson:"role_id" json:"role_id"`
	Roles   []Role               `json:"role" bson:"-"`
}
