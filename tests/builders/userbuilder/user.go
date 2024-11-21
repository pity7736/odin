package userbuilder

import (
	"github.com/google/uuid"
	"raiseexception.dev/odin/src/shared/domain/user"
)

type builder struct {
	email string
	id    string
}

func New() *builder {
	id, _ := uuid.NewV7()
	return &builder{email: "test@raiseexception.dev", id: id.String()}
}

func (b *builder) Build() *user.User {
	return user.New(b.email, b.id)
}
