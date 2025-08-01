package authen

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo *AuthenticationRepository) UpdateSession(ctx context.Context, session models.Session) error {
	_, err := repo.SessionCol.UpdateOne(ctx, bson.M{
		"user_id":   session.UserID,
		"device_id": session.DeviceID,
	}, bson.M{"$set": session})
	return err
}
