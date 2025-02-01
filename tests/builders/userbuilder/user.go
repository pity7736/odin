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

func (self *Builder) Create(repository repositories.UserRepository) *usermodel.User {
	user := self.Build()
	repository.Add(context.TODO(), user)
	return user
}

func (self *Builder) Build() *usermodel.User {
	return usermodel.New(self.id, self.email, self.password)
}
