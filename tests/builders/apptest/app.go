package apptest

import (
	accountsrepos "raiseexception.dev/odin/src/accounts/domain/repositories"
	accountsinmemory "raiseexception.dev/odin/src/accounts/infrastructure/repositories/inmemory"
	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/app"
	vaultrepos "raiseexception.dev/odin/src/vault/domain/repositories"
	vaultinmemory "raiseexception.dev/odin/src/vault/infrastructure/repositories/inmemory"
)

type Builder struct {
	sessionRepository accountsrepos.SessionRepository
	userRepository    accountsrepos.UserRepository
	chunkRepository   vaultrepos.ChunkRepository
}

func New() *Builder {
	return &Builder{
		sessionRepository: accountsinmemory.NewInMemorySessionRepository(),
		userRepository:    accountsinmemory.NewInMemoryUserRepository(),
		chunkRepository:   vaultinmemory.NewInMemoryChunkRepository(),
	}
}

func (self *Builder) WithSessionRepository(sessionRepository accountsrepos.SessionRepository) *Builder {
	self.sessionRepository = sessionRepository
	return self
}

func (self *Builder) WithUserRepository(userRepository accountsrepos.UserRepository) *Builder {
	self.userRepository = userRepository
	return self
}

func (self *Builder) WithChunkRepository(chunkRepository vaultrepos.ChunkRepository) *Builder {
	self.chunkRepository = chunkRepository
	return self
}

func (self *Builder) Build() app.Application {
	return app.NewFiberApplication(
		self.sessionRepository,
		self.userRepository,
		bcrypthasher.New(),
		self.chunkRepository,
	)
}
