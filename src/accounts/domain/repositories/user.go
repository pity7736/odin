package repositories

import (
	"context"

	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*usermodel.User, error)
	Exists(ctx context.Context, email string) (bool, error)
	Add(ctx context.Context, user *usermodel.User) error
}
