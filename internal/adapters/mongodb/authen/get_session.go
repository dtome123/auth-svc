package authen

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (svc *AuthenticationRepository) GetSession(ctx context.Context, userID, deviceID string) (*models.Session, error) {
	var session models.Session
	err := svc.SessionCol.FindOne(ctx, bson.M{
		"user_id":   userID,
		"device_id": deviceID,
	}).Decode(&session)
	return &session, err
}
