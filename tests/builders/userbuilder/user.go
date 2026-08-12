package userbuilder

import (
	"context"

	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

type Builder struct {
	email    string
	password string
}

func New() *Builder {
	return &Builder{
		email:    "test@raiseexception.dev",
		password: "some secure password",
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
	user, _ := usermodel.New(self.email, self.password)
	return user
}

func (self *Builder) Create(repository repositories.UserRepository) *usermodel.User {
	user := self.Build()
	_ = repository.Add(context.TODO(), user)
	return user
}
