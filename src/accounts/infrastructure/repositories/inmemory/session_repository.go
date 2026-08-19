package inmemory

import (
	"context"

	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type InMemorySessionRepository struct {
	sessions map[string]*sessionmodel.Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{sessions: make(map[string]*sessionmodel.Session)}
}

func (self *InMemorySessionRepository) Add(ctx context.Context, session *sessionmodel.Session) error {
	self.sessions[session.Token()] = session
	return nil
}

func (self *InMemorySessionRepository) Get(ctx context.Context, token string) (*sessionmodel.Session, error) {
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
			WithTag(odinerrors.Domain).
			Build()
	}
	return session, nil
}

func (self *InMemorySessionRepository) Save(ctx context.Context, session *sessionmodel.Session) error {
	self.sessions[session.Token()] = session
	return nil
}

func (self *InMemorySessionRepository) Delete(ctx context.Context, token string) error {
	delete(self.sessions, token)
	return nil
}
