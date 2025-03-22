package pgrepositories

import (
	"context"

	"github.com/google/uuid"

	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

type PGUserRepository struct {
	users map[string]*usermodel.User
}

func NewPGUserRepository() *PGUserRepository {
	users := make(map[string]*usermodel.User, 2)
	id, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()
	users["some@email.com"] = usermodel.New(id.String(), "some@email.com", "password")
	users["some1@email.com"] = usermodel.New(id2.String(), "some1@email.com", "password")
	return &PGUserRepository{users: users}
}

func (self *PGUserRepository) GetByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	user := self.users[email]
	return user, nil
}

func (self *PGUserRepository) Add(ctx context.Context, user *usermodel.User) error {
	self.users[user.Email()] = user
	return nil
}
