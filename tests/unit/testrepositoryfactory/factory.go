package testrepositoryfactory

import (
	"testing"

	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/tests/unit/mocks"
)

type Factory struct {
	sessionRepository *mocks.MockSessionRepository
	userRepository    *mocks.MockUserRepository
}

func New(t *testing.T) *Factory {
	return &Factory{
		sessionRepository: mocks.NewMockSessionRepository(t),
		userRepository:    mocks.NewMockUserRepository(t),
	}
}

func (self *Factory) GetSessionRepository() repositories.SessionRepository {
	return self.sessionRepository
}

func (self *Factory) GetSessionRepositoryMock() *mocks.MockSessionRepository {
	return self.sessionRepository
}

func (self *Factory) GetUserRepository() repositories.UserRepository {
	return self.userRepository
}

func (self *Factory) GetUserRepositoryMock() *mocks.MockUserRepository {
	return self.userRepository
}
