package accountsrepositoryfactory

import (
	"raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/src/accounts/infrastructure/repositories/pgrepositories"
)

type RepositoryFactory interface {
	GetUserRepository() repositories.UserRepository
	GetSessionRepository() repositories.SessionRepository
}

type repositoryFactory struct {
	userRepository    repositories.UserRepository
	sessionRepository repositories.SessionRepository
}

func New() RepositoryFactory {
	return &repositoryFactory{
		pgrepositories.NewPGUserRepository(),
		pgrepositories.NewPGSessionRepository(),
	}
}

func (self *repositoryFactory) GetUserRepository() repositories.UserRepository {
	return self.userRepository
}

func (self *repositoryFactory) GetSessionRepository() repositories.SessionRepository {
	return self.sessionRepository
}
