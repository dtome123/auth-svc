package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Assignment struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	RoleID primitive.ObjectID `bson:"role_id" json:"role_id"`
	Role   Role               `json:"role" bson:"-"`
}
