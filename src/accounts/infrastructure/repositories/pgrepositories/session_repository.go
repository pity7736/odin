package pgrepositories

import (
	"context"

	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
)

type PGSessionRepository struct {
	sessions map[string]*sessionmodel.Session
}

func NewPGSessionRepository() *PGSessionRepository {
	return &PGSessionRepository{sessions: make(map[string]*sessionmodel.Session)}
}

func (self *PGSessionRepository) Add(ctx context.Context, session *sessionmodel.Session) error {
	self.sessions[session.Token()] = session
	return nil
}

func (self *PGSessionRepository) Get(ctx context.Context, token string) (*sessionmodel.Session, error) {
	return self.sessions[token], nil
}
