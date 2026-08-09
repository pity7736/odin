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
	session := self.sessions[token]
	if session == nil || session.IsExpired() {
		return nil, nil
	}
	return session, nil
}

func (self *PGSessionRepository) Save(ctx context.Context, session *sessionmodel.Session) error {
	self.sessions[session.Token()] = session
	return nil
}

func (self *PGSessionRepository) Delete(ctx context.Context, token string) error {
	delete(self.sessions, token)
	return nil
}
