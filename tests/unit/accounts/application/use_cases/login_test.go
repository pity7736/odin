package use_cases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionstarter"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestLogin(t *testing.T) {
	t.Run("Should be able to login", func(t *testing.T) {
		builder := userbuilder.New()
		user := builder.Build()
		factory := testrepositoryfactory.New(t)
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(context.TODO(), user.Email()).Return(user, nil)
		sessionRepository := factory.GetSessionRepositoryMock()
		sessionRepository.EXPECT().Add(context.TODO(), mock.Anything).Return(nil)
		sessionStarter := sessionstarter.New(
			user.Email(),
			builder.Password(),
			factory.GetUserRepository(),
			factory.GetSessionRepository(),
			bcrypthasher.New(),
		)
		session, authenticatedUser, err := sessionStarter.Start(context.TODO())

		assert.Nil(t, err)
		assert.NotEmpty(t, session.Token())
		assert.Equal(t, user.ID(), session.UserID())
		assert.Equal(t, user, authenticatedUser)
		sessionRepository.AssertCalled(t, "Add", context.TODO(), mock.Anything)
	})

	t.Run("Should not be able to login when repository return error", func(t *testing.T) {
		builder := userbuilder.New()
		user := builder.Build()
		factory := testrepositoryfactory.New(t)
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(context.TODO(), user.Email()).Return(user, nil)
		repoErr := errors.New("error saving token to sessionRepository")
		sessionRepository := factory.GetSessionRepositoryMock()
		sessionRepository.EXPECT().Add(context.TODO(), mock.Anything).Return(repoErr)
		sessionStarter := sessionstarter.New(
			user.Email(),
			builder.Password(),
			factory.GetUserRepository(),
			factory.GetSessionRepository(),
			bcrypthasher.New(),
		)
		session, authenticatedUser, err := sessionStarter.Start(context.TODO())

		assert.Equal(t, repoErr, err)
		assert.Nil(t, session)
		assert.Nil(t, authenticatedUser)
		sessionRepository.AssertCalled(t, "Add", context.TODO(), mock.Anything)
	})

	t.Run("Should not be able to login", func(t *testing.T) {
		builder := userbuilder.New()
		user := builder.Build()
		factory := testrepositoryfactory.New(t)
		sessionRepository := factory.GetSessionRepositoryMock()
		notFoundError := odinerrors.NewErrorBuilder("user not found").
			WithExternalMessage("Usuario no encontrado").
			WithTag(odinerrors.NotFound).
			Build()
		testCases := []struct {
			name          string
			email         string
			authHash      string
			expectedUser  *usermodel.User
			expectedError error
		}{
			{
				"when auth hash is wrong",
				user.Email(),
				"wrong auth hash",
				user,
				nil,
			},
			{
				"when email is wrong",
				"wrong@test.dev",
				builder.Password(),
				nil,
				notFoundError,
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				userRepository := factory.GetUserRepositoryMock()
				userRepository.EXPECT().GetByEmail(context.TODO(), testCase.email).Return(testCase.expectedUser, testCase.expectedError)
				sessionStarter := sessionstarter.New(
					testCase.email,
					testCase.authHash,
					factory.GetUserRepository(),
					factory.GetSessionRepository(),
					bcrypthasher.New(),
				)
				session, authenticatedUser, err := sessionStarter.Start(context.TODO())

				var odinError *odinerrors.Error
				assert.True(t, errors.As(err, &odinError))
				assert.Equal(t, "Correo o contraseña incorrectos", odinError.ExternalError())
				assert.Equal(t, odinerrors.Unauthorized, odinError.Tag())
				assert.Nil(t, session)
				assert.Nil(t, authenticatedUser)
				userRepository.AssertCalled(t, "GetByEmail", context.TODO(), testCase.email)
				sessionRepository.AssertNotCalled(t, "Add", context.TODO(), mock.Anything)
			})
		}
	})

	t.Run("Should not be able to login when user repository return err", func(t *testing.T) {
		builder := userbuilder.New()
		user := builder.Build()
		factory := testrepositoryfactory.New(t)
		userRepository := factory.GetUserRepositoryMock()
		repoErr := errors.New("error getting user")
		userRepository.EXPECT().GetByEmail(context.TODO(), user.Email()).Return(nil, repoErr)
		sessionRepository := factory.GetSessionRepositoryMock()
		sessionStarter := sessionstarter.New(
			user.Email(),
			builder.Password(),
			factory.GetUserRepository(),
			factory.GetSessionRepository(),
			bcrypthasher.New(),
		)
		session, authenticatedUser, err := sessionStarter.Start(context.TODO())

		assert.Equal(t, repoErr, err)
		assert.Nil(t, session)
		assert.Nil(t, authenticatedUser)
		sessionRepository.AssertNotCalled(t, "Add", context.TODO(), mock.Anything)
	})
}
