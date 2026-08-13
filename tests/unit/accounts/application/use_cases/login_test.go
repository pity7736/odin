package use_cases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionstarter"
	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/tests/builders/userbuilder"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestSessionStarterShould(t *testing.T) {
	t.Run("start session successfully", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		builder := userbuilder.New()
		user := builder.Build()
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(context.TODO(), user.Email()).Return(user, nil)
		sessionRepository := factory.GetSessionRepositoryMock()
		sessionRepository.EXPECT().Add(context.TODO(), mock.Anything).Return(nil)
		starter := sessionstarter.New(user.Email(), builder.Password(), factory.GetUserRepository(), factory.GetSessionRepository(), bcrypthasher.New())
		session, err := starter.Start(context.TODO())
		assert.Nil(t, err)
		assert.NotEmpty(t, session.Token())
		assert.Equal(t, user.ID(), session.UserID())
	})
	t.Run("return error when password is wrong", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		builder := userbuilder.New()
		user := builder.Build()
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(context.TODO(), user.Email()).Return(user, nil)
		starter := sessionstarter.New(user.Email(), "wrong password", factory.GetUserRepository(), factory.GetSessionRepository(), bcrypthasher.New())
		session, err := starter.Start(context.TODO())
		assert.Nil(t, session)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, "Correo o contraseña incorrectos", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})
	t.Run("return error when user is not found", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(context.TODO(), "unknown@example.com").Return(nil, nil)
		starter := sessionstarter.New("unknown@example.com", "password", factory.GetUserRepository(), factory.GetSessionRepository(), bcrypthasher.New())
		session, err := starter.Start(context.TODO())
		assert.Nil(t, session)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, "Correo o contraseña incorrectos", odinError.ExternalError())
	})
	t.Run("propagate user repository error", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		repoError := errors.New("database failure")
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(context.TODO(), "test@example.com").Return(nil, repoError)
		starter := sessionstarter.New("test@example.com", "password", factory.GetUserRepository(), factory.GetSessionRepository(), bcrypthasher.New())
		session, err := starter.Start(context.TODO())
		assert.Nil(t, session)
		assert.Equal(t, repoError, err)
	})
	t.Run("propagate session repository error", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		builder := userbuilder.New()
		user := builder.Build()
		repoError := errors.New("session store failure")
		userRepository := factory.GetUserRepositoryMock()
		userRepository.EXPECT().GetByEmail(context.TODO(), user.Email()).Return(user, nil)
		sessionRepository := factory.GetSessionRepositoryMock()
		sessionRepository.EXPECT().Add(context.TODO(), mock.Anything).Return(repoError)
		starter := sessionstarter.New(user.Email(), builder.Password(), factory.GetUserRepository(), factory.GetSessionRepository(), bcrypthasher.New())
		session, err := starter.Start(context.TODO())
		assert.Nil(t, session)
		assert.Equal(t, repoError, err)
	})
}
