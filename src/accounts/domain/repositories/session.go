package repositories

import "context"

type SessionRepository interface {
	Add(ctx context.Context, token string) error
}
