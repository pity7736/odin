package userbuilder

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"raiseexception.dev/odin/src/accounts/domain/keyparams"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
)

const DefaultPassword = "some secure password"
const DefaultEncryptedMasterKey = "encrypted-master-key-base64"
const DefaultAlgorithm = "argon2id"
const DefaultIterations = 3
const DefaultMemory = 65536
const DefaultParallelism = 4
const DefaultSalt = "default-salt-base64"

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
	params, _ := keyparams.New(DefaultAlgorithm, DefaultIterations, DefaultMemory, DefaultParallelism, DefaultSalt)
	user, _ := usermodel.New(self.email, string(hash), DefaultEncryptedMasterKey, params)
	return user
}

func (self *Builder) Create(repository repositories.UserRepository) *usermodel.User {
	user := self.Build()
	_ = repository.Add(context.TODO(), user)
	return user
}
