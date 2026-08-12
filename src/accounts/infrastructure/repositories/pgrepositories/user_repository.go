package pgrepositories

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

type PGUserRepository struct {
	users map[string]*usermodel.User
}

func NewPGUserRepository() *PGUserRepository {
	users := make(map[string]*usermodel.User, 2)
	user1, _ := usermodel.New("some@email.com", hashPassword("password"))
	user2, _ := usermodel.New("some1@email.com", hashPassword("password"))
	users["some@email.com"] = user1
	users["some1@email.com"] = user2
	return &PGUserRepository{users: users}
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func (self *PGUserRepository) GetByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	user := self.users[email]
	return user, nil
}

func (self *PGUserRepository) Add(ctx context.Context, user *usermodel.User) error {
	self.users[user.Email()] = user
	return nil
}
