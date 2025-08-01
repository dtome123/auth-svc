package authen

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthenticationRepository struct {
	SessionCol *mongo.Collection
}

func NewAuthenticationRepository(db *mongo.Database) *AuthenticationRepository {

	sessionCol := db.Collection("sessions")

	indexingSessionCol(context.TODO(), sessionCol)

	return &AuthenticationRepository{
		SessionCol: sessionCol,
	}
}
