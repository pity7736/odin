package pgrepositories

import (
	"context"
	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
)

type PGSessionRepository struct {
	sessions []*sessionmodel.Session
}

func NewPGSessionRepository() *PGSessionRepository {
	return &PGSessionRepository{sessions: make([]*sessionmodel.Session, 0)}
}

func (self *PGSessionRepository) Add(ctx context.Context, session *sessionmodel.Session) error {
	self.sessions = append(self.sessions, session)
	return nil
}
