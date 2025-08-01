package authen

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (svc *AuthenticationRepository) GetSession(ctx context.Context, userID, deviceID string) (*models.Session, error) {
	var session models.Session
	err := svc.SessionCol.FindOne(ctx, bson.M{
		"user_id":   userID,
		"device_id": deviceID,
	}).Decode(&session)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &session, nil
}
