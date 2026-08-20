package testrepositoryfactory

import (
	"testing"

	accountsrepos "raiseexception.dev/odin/src/accounts/domain/repositories"
	vaultrepos "raiseexception.dev/odin/src/vault/domain/repositories"
	"raiseexception.dev/odin/tests/unit/mocks"
)

type Factory struct {
	sessionRepository *mocks.MockSessionRepository
	userRepository    *mocks.MockUserRepository
	chunkRepository   *mocks.MockChunkRepository
}

func New(t *testing.T) *Factory {
	return &Factory{
		sessionRepository: mocks.NewMockSessionRepository(t),
		userRepository:    mocks.NewMockUserRepository(t),
		chunkRepository:   mocks.NewMockChunkRepository(t),
	}
}

func (self *Factory) GetSessionRepository() accountsrepos.SessionRepository {
	return self.sessionRepository
}

func (self *Factory) GetSessionRepositoryMock() *mocks.MockSessionRepository {
	return self.sessionRepository
}

func (self *Factory) GetUserRepository() accountsrepos.UserRepository {
	return self.userRepository
}

func (self *Factory) GetUserRepositoryMock() *mocks.MockUserRepository {
	return self.userRepository
}

func (self *Factory) GetChunkRepository() vaultrepos.ChunkRepository {
	return self.chunkRepository
}

func (self *Factory) GetChunkRepositoryMock() *mocks.MockChunkRepository {
	return self.chunkRepository
}
