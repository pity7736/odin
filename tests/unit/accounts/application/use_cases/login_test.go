package login_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionstarter"
	"raiseexception.dev/odin/src/accounts/domain/usermodel"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestLogin(t *testing.T) {
	t.Run("Should be able to login", func(t *testing.T) {
		user := userbuilder.New().Build()
		factory := testrepositoryfactory.New(t)
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(context.TODO(), user.Email()).Return(user, nil)
		sessionRepository := factory.GetSessionRepositoryMock()
		sessionRepository.EXPECT().Add(context.TODO(), mock.Anything).Return(nil)
		sessionStarter := sessionstarter.New(
			user.Email(),
			user.Password(),
			factory,
		)
		session, err := sessionStarter.Start(context.TODO())

		assert.Nil(t, err)
		assert.NotEmpty(t, session.Token())
		sessionRepository.AssertCalled(t, "Add", context.TODO(), mock.Anything)
	})

	t.Run("Should not be able to login when repository return error", func(t *testing.T) {
		user := userbuilder.New().Build()
		factory := testrepositoryfactory.New(t)
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(context.TODO(), user.Email()).Return(user, nil)
		repoErr := errors.New("error saving token to sessionRepository")
		sessionRepository := factory.GetSessionRepositoryMock()
		sessionRepository.EXPECT().Add(context.TODO(), mock.Anything).Return(repoErr)
		sessionStarter := sessionstarter.New(
			user.Email(),
			user.Password(),
			factory,
		)
		session, err := sessionStarter.Start(context.TODO())

		assert.Equal(t, repoErr, err)
		assert.Nil(t, session)
		sessionRepository.AssertCalled(t, "Add", context.TODO(), mock.Anything)
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
				userRepository.EXPECT().GetByEmail(context.TODO(), testCase.email).Return(testCase.expectedUser, nil)
				repoErr := errors.New("email or password are wrong")
				sessionStarter := sessionstarter.New(
					testCase.email,
					testCase.password,
					factory,
				)
				session, err := sessionStarter.Start(context.TODO())

				assert.Equal(t, repoErr, err)
				assert.Nil(t, session)
				userRepository.AssertCalled(t, "GetByEmail", context.TODO(), testCase.email)
				sessionRepository.AssertNotCalled(t, "Add", context.TODO(), mock.Anything)
			})
		}
	})

	t.Run("Should not be able to login when user repository return err", func(t *testing.T) {
		user := userbuilder.New().Build()
		factory := testrepositoryfactory.New(t)
		userRepository := factory.GetUserRepositoryMock()
		repoErr := errors.New("error getting user")
		userRepository.EXPECT().GetByEmail(context.TODO(), user.Email()).Return(nil, repoErr)
		sessionRepository := factory.GetSessionRepositoryMock()
		sessionStarter := sessionstarter.New(
			user.Email(),
			user.Password(),
			factory,
		)
		session, err := sessionStarter.Start(context.TODO())

		assert.Equal(t, repoErr, err)
		assert.Nil(t, session)
		sessionRepository.AssertNotCalled(t, "Add", context.TODO(), mock.Anything)
	})
}
