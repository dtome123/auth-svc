package authen

import (
	"auth-svc/internal/models"

	mongodb "github.com/dtome123/go-mongo-generic"
)

type AuthenticationRepository struct {
	SessionCol mongodb.Collection[models.Session]
}

func NewAuthenticationRepository(db *mongodb.Database) *AuthenticationRepository {

	sessionCol := mongodb.NewCollection[models.Session](db)
	sessionCol.EnsureIndexes(GetSessionIndexes())

	return &AuthenticationRepository{
		SessionCol: sessionCol,
	}
}
