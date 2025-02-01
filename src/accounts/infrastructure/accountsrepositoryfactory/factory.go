package accountsrepositoryfactory

import (
	"raiseexception.dev/odin/src/accounts/domain/repositories"
)

type AccountsRepositoryFactory interface {
	GetSessionRepository() repositories.SessionRepository
	GetUserRepository() repositories.UserRepository
}
