package api

import "context"

type Handler interface {
	Handle(ctx context.Context) error
}
