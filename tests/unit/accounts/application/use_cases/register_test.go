package use_cases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"raiseexception.dev/odin/src/accounts/application/use_cases/userregistrar"
	"raiseexception.dev/odin/src/accounts/domain/keyparams"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func newKeyParams() keyparams.KeyParams {
	params, _ := keyparams.New(
		userbuilder.DefaultAlgorithm,
		userbuilder.DefaultIterations,
		userbuilder.DefaultMemory,
		userbuilder.DefaultParallelism,
		userbuilder.DefaultSalt,
	)
	return params
}

func TestRegisterShould(t *testing.T) {
	t.Run("register a new user when the email is available", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		email := "new@example.com"
		authHash := "client-derived-auth-hash"
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(context.TODO(), email).Return(nil, nil)
		var storedUser *usermodel.User
		userRepository.EXPECT().Add(context.TODO(), mock.Anything).Run(func(ctx context.Context, user *usermodel.User) {
			storedUser = user
		}).Return(nil)
		registrar := userregistrar.New(
			email,
			authHash,
			userbuilder.DefaultEncryptedMasterKey,
			newKeyParams(),
			factory.GetUserRepository(),
			bcrypthasher.New(),
		)
		user, err := registrar.Register(context.TODO())

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email())
		assert.Equal(t, userbuilder.DefaultEncryptedMasterKey, user.EncryptedMasterKey())
		assert.Equal(t, userbuilder.DefaultAlgorithm, user.KeyParams().Algorithm())
		assert.True(t, bcrypthasher.New().Compare(user.AuthHashDigest(), authHash))
		assert.Equal(t, user, storedUser)
		userRepository.AssertCalled(t, "Add", context.TODO(), mock.Anything)
	})

	t.Run("reject registration when the email already exists", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		existingUser := userbuilder.New().WithEmail("taken@example.com").Build()
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(context.TODO(), existingUser.Email()).Return(existingUser, nil)
		registrar := userregistrar.New(
			existingUser.Email(),
			"client-derived-auth-hash",
			userbuilder.DefaultEncryptedMasterKey,
			newKeyParams(),
			factory.GetUserRepository(),
			bcrypthasher.New(),
		)
		user, err := registrar.Register(context.TODO())

		assert.Nil(t, user)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, "El correo ya está registrado", odinError.ExternalError())
		assert.Equal(t, odinerrors.AlreadyExists, odinError.Tag())
		userRepository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("propagate the error when looking up the email fails", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		email := "new@example.com"
		lookupError := errors.New("error getting user")
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(context.TODO(), email).Return(nil, lookupError)
		registrar := userregistrar.New(
			email,
			"client-derived-auth-hash",
			userbuilder.DefaultEncryptedMasterKey,
			newKeyParams(),
			factory.GetUserRepository(),
			bcrypthasher.New(),
		)
		user, err := registrar.Register(context.TODO())

		assert.Nil(t, user)
		assert.Equal(t, lookupError, err)
		userRepository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("propagate the error when persisting the user fails", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		email := "new@example.com"
		persistError := errors.New("error saving user")
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(context.TODO(), email).Return(nil, nil)
		userRepository.EXPECT().Add(context.TODO(), mock.Anything).Return(persistError)
		registrar := userregistrar.New(
			email,
			"client-derived-auth-hash",
			userbuilder.DefaultEncryptedMasterKey,
			newKeyParams(),
			factory.GetUserRepository(),
			bcrypthasher.New(),
		)
		user, err := registrar.Register(context.TODO())

		assert.Nil(t, user)
		assert.Equal(t, persistError, err)
	})
}
