package inmemory

import (
	"context"

	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
)

type InMemoryChunkRepository struct {
	chunks map[string]map[string]*chunkmodel.EncryptedChunk
}

func NewInMemoryChunkRepository() *InMemoryChunkRepository {
	return &InMemoryChunkRepository{chunks: make(map[string]map[string]*chunkmodel.EncryptedChunk)}
}

func (self *InMemoryChunkRepository) GetByID(ctx context.Context, ownerID, id string) (*chunkmodel.EncryptedChunk, error) {
	ownedChunks, ok := self.chunks[ownerID]
	if !ok {
		return nil, nil
	}
	return ownedChunks[id], nil
}

func (self *InMemoryChunkRepository) Add(ctx context.Context, chunk *chunkmodel.EncryptedChunk) error {
	ownedChunks, ok := self.chunks[chunk.OwnerID()]
	if !ok {
		ownedChunks = make(map[string]*chunkmodel.EncryptedChunk)
		self.chunks[chunk.OwnerID()] = ownedChunks
	}
	ownedChunks[chunk.ID()] = chunk
	return nil
}
