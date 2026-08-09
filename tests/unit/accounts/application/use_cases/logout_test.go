package use_cases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/accounts/application/use_cases/sessionterminator"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestSessionTerminatorShould(t *testing.T) {
	t.Run("terminate session successfully", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		sessionRepository := factory.GetSessionRepositoryMock()
		sessionRepository.EXPECT().Delete(context.TODO(), "session-token").Return(nil)
		terminator := sessionterminator.New(factory.GetSessionRepository())
		err := terminator.Terminate(context.TODO(), "session-token")
		assert.Nil(t, err)
	})
	t.Run("propagate session repository error", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		repoError := errors.New("delete failure")
		sessionRepository := factory.GetSessionRepositoryMock()
		sessionRepository.EXPECT().Delete(context.TODO(), "session-token").Return(repoError)
		terminator := sessionterminator.New(factory.GetSessionRepository())
		err := terminator.Terminate(context.TODO(), "session-token")
		assert.Equal(t, repoError, err)
	})
}
