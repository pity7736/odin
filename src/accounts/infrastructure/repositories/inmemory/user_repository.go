package inmemory

import (
	"context"

	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

type InMemoryUserRepository struct {
	users map[string]*usermodel.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[string]*usermodel.User)}
}

func (self *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	user := self.users[email]
	return user, nil
}

func (self *InMemoryUserRepository) Add(ctx context.Context, user *usermodel.User) error {
	self.users[user.Email()] = user
	return nil
}
