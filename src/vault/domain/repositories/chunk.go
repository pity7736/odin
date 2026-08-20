package repositories

import (
	"context"

	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
)

type ChunkRepository interface {
	Exists(ctx context.Context, ownerID, id string) (bool, error)
	Add(ctx context.Context, chunk *chunkmodel.EncryptedChunk) error
	Get(ctx context.Context, ownerID, id string) (*chunkmodel.EncryptedChunk, error)
}
