package userbuilder

import (
	"github.com/google/uuid"
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

func (self *Builder) Build() *usermodel.User {
	return usermodel.New(self.email, self.id, self.password)
}
