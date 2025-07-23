package authen

import (
	"auth-svc/internal/models"
	"context"
)


func (repo *AuthenticationRepository) CreateSession(ctx context.Context, session models.Session) error {
	_, err := repo.SessionCol.InsertOne(ctx, session)
	return err
}
