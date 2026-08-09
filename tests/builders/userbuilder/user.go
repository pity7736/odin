package userbuilder

import (
	"context"

	"github.com/google/uuid"

	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

type Builder struct {
	email    string
	id       string
	password string
}

func New() *Builder {
	id, _ := uuid.NewV7()
	return &Builder{
		email:    "test@raiseexception.dev",
		id:       id.String(),
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
	user, _ := usermodel.NewWithPlainPassword(self.id, self.email, self.password)
	return user
}

func (self *Builder) Create(repository repositories.UserRepository) *usermodel.User {
	user := self.Build()
	_ = repository.Add(context.TODO(), user)
	return user
}
