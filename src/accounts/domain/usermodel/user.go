package usermodel

import (
	"github.com/google/uuid"
	"raiseexception.dev/odin/src/accounts/domain/keyparams"
)

type User struct {
	id                 string
	email              string
	authHashDigest     string
	encryptedMasterKey string
	keyParams          keyparams.KeyParams
}

func New(email, authHashDigest, encryptedMasterKey string, keyParams keyparams.KeyParams) (*User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &User{
		id:                 id.String(),
		email:              email,
		authHashDigest:     authHashDigest,
		encryptedMasterKey: encryptedMasterKey,
		keyParams:          keyParams,
	}, nil
}

func (self *User) ID() string {
	return self.id
}

func (self *User) Email() string {
	return self.email
}

func (self *User) AuthHashDigest() string {
	return self.authHashDigest
}

func (self *User) EncryptedMasterKey() string {
	return self.encryptedMasterKey
}

func (self *User) KeyParams() keyparams.KeyParams {
	return self.keyParams
}
