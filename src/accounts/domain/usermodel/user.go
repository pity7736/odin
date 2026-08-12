package usermodel

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id             string
	email          string
	hashedPassword string
}

func New(email, password string) (*User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{id: id.String(), email: email, hashedPassword: string(hash)}, nil
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

func (self *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(self.hashedPassword), []byte(password)) == nil
}
