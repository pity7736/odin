package repositories

import (
	"context"

	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
)

type ChunkRepository interface {
	GetByID(ctx context.Context, ownerID, id string) (*chunkmodel.EncryptedChunk, error)
	Add(ctx context.Context, chunk *chunkmodel.EncryptedChunk) error
}
