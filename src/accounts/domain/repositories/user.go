package repositories

import "raiseexception.dev/odin/src/accounts/domain/usermodel"

type UserRepository interface {
	GetByEmail(email string) (*usermodel.User, error)
}
