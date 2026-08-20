package inmemory

import (
	"context"
	"sort"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
)

type InMemoryChunkRepository struct {
	chunks map[string]map[string]*chunkmodel.EncryptedChunk
}

func NewInMemoryChunkRepository() *InMemoryChunkRepository {
	return &InMemoryChunkRepository{chunks: make(map[string]map[string]*chunkmodel.EncryptedChunk)}
}

func (self *InMemoryChunkRepository) Exists(ctx context.Context, ownerID, id string) (bool, error) {
	ownedChunks, ok := self.chunks[ownerID]
	if !ok {
		return false, nil
	}
	_, ok = ownedChunks[id]
	return ok, nil
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

func (self *InMemoryChunkRepository) Get(ctx context.Context, ownerID, id string) (*chunkmodel.EncryptedChunk, error) {
	ownedChunks, ok := self.chunks[ownerID]
	if !ok {
		return nil, self.notFound()
	}
	chunk, ok := ownedChunks[id]
	if !ok {
		return nil, self.notFound()
	}
	return chunk, nil
}

func (self *InMemoryChunkRepository) GetAll(ctx context.Context, ownerID string) ([]*chunkmodel.EncryptedChunk, error) {
	ownedChunks := self.chunks[ownerID]
	chunks := make([]*chunkmodel.EncryptedChunk, 0, len(ownedChunks))
	for _, chunk := range ownedChunks {
		chunks = append(chunks, chunk)
	}
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].ID() > chunks[j].ID()
	})
	return chunks, nil
}

func (self *InMemoryChunkRepository) notFound() error {
	return odinerrors.NewErrorBuilder("chunk not found").
		WithExternalMessage("El elemento no existe").
		WithTag(odinerrors.NotFound).
		Build()
}
