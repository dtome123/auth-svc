package authen

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *AuthenticationRepository) UpsertSession(ctx context.Context, session models.Session) error {
	_, err := repo.SessionCol.UpdateOne(ctx, bson.M{
		"user_id":   session.UserID,
		"device_id": session.DeviceID,
	}, bson.M{"$set": session}, options.Update().SetUpsert(true))
	return err
}
