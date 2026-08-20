package chunkgetter

import (
	"context"

	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
	"raiseexception.dev/odin/src/vault/domain/repositories"
)

type ChunkGetter struct {
	id              string
	ownerID         string
	chunkRepository repositories.ChunkRepository
}

func New(id, ownerID string, chunkRepository repositories.ChunkRepository) ChunkGetter {
	return ChunkGetter{
		id:              id,
		ownerID:         ownerID,
		chunkRepository: chunkRepository,
	}
}

func (self ChunkGetter) Get(ctx context.Context) (*chunkmodel.EncryptedChunk, error) {
	return self.chunkRepository.Get(ctx, self.ownerID, self.id)
}
