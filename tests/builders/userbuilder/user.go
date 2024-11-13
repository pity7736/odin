package userbuilder

import "raiseexception.dev/odin/src/shared/domain/user"

type builder struct {
	email string
}

func New() *builder {
	return &builder{email: "test@raiseexception.dev"}
}

func (b *builder) Build() *user.User {
	return user.New(b.email)
}
