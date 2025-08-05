package authen

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *AuthenticationRepository) UpdateSession(ctx context.Context, session models.Session) error {
	err := repo.SessionCol.UpdateSetOne(ctx, bson.M{
		"user_id":   session.UserID,
		"device_id": session.DeviceID,
	}, bson.M{"$set": session}, options.Update().SetHint(IdxSessionUserIdDeviceId))
	return err
}
