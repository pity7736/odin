package pgrepositories

import (
	"context"
)

type PGSessionRepository struct {
	tokens []string
}

func NewPGSessionRepository() *PGSessionRepository {
	return &PGSessionRepository{tokens: make([]string, 0)}
}

func (self *PGSessionRepository) Add(ctx context.Context, token string) error {
	self.tokens = append(self.tokens, token)
	return nil
}
