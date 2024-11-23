package login_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionstarter"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Run("Should be able to login", func(t *testing.T) {
		user := userbuilder.New().Build()
		factory := testrepositoryfactory.New(t)
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(user.Email()).Return(user, nil)
		sessionRepository := factory.GetSessionRepositoryMock()
		sessionRepository.EXPECT().Add(mock.Anything).Return(nil)
		sessionStarter := sessionstarter.New(
			user.Email(),
			user.Password(),
			factory,
		)
		token, err := sessionStarter.Start()

		assert.Nil(t, err)
		assert.NotEmpty(t, token)
		sessionRepository.AssertCalled(t, "Add", mock.Anything)
	})

	t.Run("Should not be able to login when repository return error", func(t *testing.T) {
		user := userbuilder.New().Build()
		factory := testrepositoryfactory.New(t)
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(user.Email()).Return(user, nil)
		repoErr := errors.New("error saving token to sessionRepository")
		sessionRepository := factory.GetSessionRepositoryMock()
		sessionRepository.EXPECT().Add(mock.Anything).Return(repoErr)
		sessionStarter := sessionstarter.New(
			user.Email(),
			user.Password(),
			factory,
		)
		token, err := sessionStarter.Start()

		assert.Equal(t, repoErr, err)
		assert.Empty(t, token)
		sessionRepository.AssertCalled(t, "Add", mock.Anything)
	})

	t.Run("Should not be able to login", func(t *testing.T) {
		user := userbuilder.New().Build()
		factory := testrepositoryfactory.New(t)
		sessionRepository := factory.GetSessionRepositoryMock()
		testCases := []struct {
			name         string
			email        string
			password     string
			expectedUser *usermodel.User
		}{
			{
				"when password is wrong",
				user.Email(),
				"wrong password",
				user,
			},
			{
				"when email is wrong",
				"wrong@test.dev",
				user.Password(),
				nil,
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				userRepository := factory.GetUserRepositoryMock()
				userRepository.EXPECT().GetByEmail(testCase.email).Return(testCase.expectedUser, nil)
				repoErr := errors.New("email or password are wrong")
				sessionStarter := sessionstarter.New(
					testCase.email,
					testCase.password,
					factory,
				)
				token, err := sessionStarter.Start()

				assert.Equal(t, repoErr, err)
				assert.Empty(t, token)
				userRepository.AssertCalled(t, "GetByEmail", testCase.email)
				sessionRepository.AssertNotCalled(t, "Add", mock.Anything)
			})
		}
	})

	t.Run("Should not be able to login when user repository return err", func(t *testing.T) {
		user := userbuilder.New().Build()
		factory := testrepositoryfactory.New(t)
		userRepository := factory.GetUserRepositoryMock()
		repoErr := errors.New("error getting user")
		userRepository.EXPECT().GetByEmail(user.Email()).Return(nil, repoErr)
		sessionRepository := factory.GetSessionRepositoryMock()
		sessionStarter := sessionstarter.New(
			user.Email(),
			user.Password(),
			factory,
		)
		token, err := sessionStarter.Start()

		assert.Equal(t, repoErr, err)
		assert.Empty(t, token)
		sessionRepository.AssertNotCalled(t, "Add", mock.Anything)
	})
}
