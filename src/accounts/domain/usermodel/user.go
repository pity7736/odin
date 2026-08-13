package usermodel

import "github.com/google/uuid"

type User struct {
	id             string
	email          string
	hashedPassword string
}

func New(email, hashedPassword string) (*User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &User{id: id.String(), email: email, hashedPassword: hashedPassword}, nil
}

func (self *User) ID() string {
	return self.id
}

func (self *User) Email() string {
	return self.email
}

func (self *User) HashedPassword() string {
	return self.hashedPassword
}
