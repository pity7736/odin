package chunkcreator

import (
	"context"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
	"raiseexception.dev/odin/src/vault/domain/repositories"
)

type ChunkCreator struct {
	id              string
	ownerID         string
	content         string
	chunkRepository repositories.ChunkRepository
}

func New(id, ownerID, content string, chunkRepository repositories.ChunkRepository) ChunkCreator {
	return ChunkCreator{
		id:              id,
		ownerID:         ownerID,
		content:         content,
		chunkRepository: chunkRepository,
	}
}

func (self ChunkCreator) Create(ctx context.Context) (*chunkmodel.EncryptedChunk, error) {
	exists, err := self.chunkRepository.Exists(ctx, self.ownerID, self.id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, odinerrors.NewErrorBuilder("chunk already exists").
			WithExternalMessage("El elemento ya existe").
			WithTag(odinerrors.AlreadyExists).
			Build()
	}
	return self.create(ctx)
}

func (self ChunkCreator) create(ctx context.Context) (*chunkmodel.EncryptedChunk, error) {
	chunk, err := chunkmodel.New(self.id, self.ownerID, self.content)
	if err != nil {
		return nil, err
	}
	if err := self.chunkRepository.Add(ctx, chunk); err != nil {
		return nil, err
	}
	return chunk, nil
}
