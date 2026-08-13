package sessionterminator

import (
	"context"

	"raiseexception.dev/odin/src/accounts/domain/repositories"
)

type SessionTerminator struct {
	sessionRepository repositories.SessionRepository
}

func New(sessionRepository repositories.SessionRepository) SessionTerminator {
	return SessionTerminator{sessionRepository: sessionRepository}
}

func (self SessionTerminator) Terminate(ctx context.Context, token string) error {
	return self.sessionRepository.Delete(ctx, token)
}
