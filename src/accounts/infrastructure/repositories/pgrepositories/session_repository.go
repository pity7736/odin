package pgrepositories

import (
	"context"
<<<<<<< HEAD
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
=======
)

type PGSessionRepository struct {
	tokens []string
}

func NewPGSessionRepository() *PGSessionRepository {
	return &PGSessionRepository{tokens: make([]string, 0)}
}

func (self *PGSessionRepository) Add(ctx context.Context, token string) error {
	self.tokens = append(self.tokens, token)
>>>>>>> 76b538cbcd0230565265d556fc314109638d19bb
	return nil
}
