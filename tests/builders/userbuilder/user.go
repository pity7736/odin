package userbuilder

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

const DefaultPassword = "some secure password"

type Builder struct {
	email    string
	password string
}

func New() *Builder {
	return &Builder{
		email:    "test@raiseexception.dev",
		password: DefaultPassword,
	}
}

func (self *Builder) WithEmail(email string) *Builder {
	self.email = email
	return self
}

func (self *Builder) Password() string {
	return self.password
}

func (self *Builder) Build() *usermodel.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(self.password), bcrypt.DefaultCost)
	user, _ := usermodel.New(self.email, string(hash))
	return user
}

func (self *Builder) Create(repository repositories.UserRepository) *usermodel.User {
	user := self.Build()
	_ = repository.Add(context.TODO(), user)
	return user
}
