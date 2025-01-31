package pgrepositories

import (
	"context"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

type PGUserRepository struct {
	users map[string]*usermodel.User
}

func NewPGUserRepository() *PGUserRepository {
	return &PGUserRepository{users: make(map[string]*usermodel.User)}
}

func (self *PGUserRepository) GetByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	user := self.users[email]
	return user, nil
}

func (self *PGUserRepository) Add(ctx context.Context, user *usermodel.User) error {
	self.users[user.Email()] = user
	return nil
}
