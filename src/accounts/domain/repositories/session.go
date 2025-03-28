package repositories

import (
	"context"

	"raiseexception.dev/odin/src/accounts/domain/sessionmodel"
)

type SessionRepository interface {
	Add(ctx context.Context, session *sessionmodel.Session) error
	Get(ctx context.Context, token string) (*sessionmodel.Session, error)
}
