package inmemory

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"raiseexception.dev/odin/src/accounts/domain/keyparams"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

type InMemoryUserRepository struct {
	users map[string]*usermodel.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	users := make(map[string]*usermodel.User, 2)
	params, _ := keyparams.New("argon2id", 3, 65536, 4, "seed-salt-base64")
	user1, _ := usermodel.New("some@email.com", hashPassword("password"), "encrypted-master-key-1", params)
	user2, _ := usermodel.New("some1@email.com", hashPassword("password"), "encrypted-master-key-2", params)
	users["some@email.com"] = user1
	users["some1@email.com"] = user2
	return &InMemoryUserRepository{users: users}
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func (self *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	user := self.users[email]
	return user, nil
}

func (self *InMemoryUserRepository) Add(ctx context.Context, user *usermodel.User) error {
	self.users[user.Email()] = user
	return nil
}
