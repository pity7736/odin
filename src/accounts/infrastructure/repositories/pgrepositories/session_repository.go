package pgrepositories

import (
	"context"

	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
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
	if session == nil {
		return nil, odinerrors.NewErrorBuilder("session not found").
			WithExternalMessage("Sesión no encontrada").
			WithTag(odinerrors.NotFound).
			Build()
	}
	if session.IsExpired() {
		delete(self.sessions, token)
		return nil, odinerrors.NewErrorBuilder("session expired").
			WithExternalMessage("Sesión expirada").
			WithTag(odinerrors.DOMAIN).
			Build()
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
