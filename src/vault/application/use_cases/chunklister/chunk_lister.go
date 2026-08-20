package chunklister

import (
	"context"

	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
	"raiseexception.dev/odin/src/vault/domain/repositories"
)

type ChunkLister struct {
	ownerID         string
	chunkRepository repositories.ChunkRepository
}

func New(ownerID string, chunkRepository repositories.ChunkRepository) ChunkLister {
	return ChunkLister{
		ownerID:         ownerID,
		chunkRepository: chunkRepository,
	}
}

func (self ChunkLister) List(ctx context.Context) ([]*chunkmodel.EncryptedChunk, error) {
	return self.chunkRepository.GetAll(ctx, self.ownerID)
}
