package inmemory

import (
	"context"

	"raiseexception.dev/odin/src/accounts/domain/usermodel"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type InMemoryUserRepository struct {
	users map[string]*usermodel.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[string]*usermodel.User)}
}

func (self *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	user, ok := self.users[email]
	if !ok {
		return nil, odinerrors.NewErrorBuilder("user not found").
			WithExternalMessage("Usuario no encontrado").
			WithTag(odinerrors.NotFound).
			Build()
	}
	return user, nil
}

func (self *InMemoryUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	_, ok := self.users[email]
	return ok, nil
}

func (self *InMemoryUserRepository) Add(ctx context.Context, user *usermodel.User) error {
	self.users[user.Email()] = user
	return nil
}
