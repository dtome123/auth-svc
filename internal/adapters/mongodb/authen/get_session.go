package authen

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (svc *AuthenticationRepository) GetSession(ctx context.Context, userID, deviceID string) (*models.Session, error) {

	session, err := svc.SessionCol.FindOne(ctx, bson.M{
		"user_id":   userID,
		"device_id": deviceID,
	}, options.FindOne().SetHint(IdxSessionUserIdDeviceId))

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return session, nil
}
