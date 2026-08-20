package userregistrar

import (
	"context"

	"raiseexception.dev/odin/src/accounts/application/authhasher"
	"raiseexception.dev/odin/src/accounts/domain/keyparams"
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type UserRegistrar struct {
	keyParams          keyparams.KeyParams
	email              string
	authHash           string
	encryptedMasterKey string
	userRepository     repositories.UserRepository
	authHasher         authhasher.AuthHasher
}

func New(
	email, authHash, encryptedMasterKey string,
	keyParams keyparams.KeyParams,
	userRepository repositories.UserRepository,
	authHasher authhasher.AuthHasher,
) UserRegistrar {

	return UserRegistrar{
		keyParams:          keyParams,
		email:              email,
		authHash:           authHash,
		encryptedMasterKey: encryptedMasterKey,
		userRepository:     userRepository,
		authHasher:         authHasher,
	}
}

func (self UserRegistrar) Register(ctx context.Context) (*usermodel.User, error) {
	exists, err := self.userRepository.Exists(ctx, self.email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, odinerrors.NewErrorBuilder("email already registered").
			WithExternalMessage("El correo ya está registrado").
			WithTag(odinerrors.AlreadyExists).
			Build()
	}
	return self.register(ctx)
}

func (self UserRegistrar) register(ctx context.Context) (*usermodel.User, error) {
	digest, err := self.authHasher.Hash(self.authHash)
	if err != nil {
		return nil, err
	}
	user, err := usermodel.New(self.email, digest, self.encryptedMasterKey, self.keyParams)
	if err != nil {
		return nil, err
	}
	err = self.userRepository.Add(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
