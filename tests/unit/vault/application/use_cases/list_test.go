package use_cases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/vault/application/use_cases/chunklister"
	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestListChunksShould(t *testing.T) {
	t.Run("return the owner's chunks unchanged", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		ownerID := "owner-id"
		firstChunk, _ := chunkmodel.New(uuid.NewString(), ownerID, "first-content")
		secondChunk, _ := chunkmodel.New(uuid.NewString(), ownerID, "second-content")
		storedChunks := []*chunkmodel.EncryptedChunk{firstChunk, secondChunk}
		chunkRepository := factory.GetChunkRepositoryMock()
		chunkRepository.EXPECT().GetAll(context.TODO(), ownerID).Return(storedChunks, nil)
		lister := chunklister.New(ownerID, factory.GetChunkRepository())
		chunks, err := lister.List(context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, storedChunks, chunks)
	})

	t.Run("return an empty slice and no error when the owner has no chunks", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		ownerID := "owner-id"
		emptyChunks := []*chunkmodel.EncryptedChunk{}
		chunkRepository := factory.GetChunkRepositoryMock()
		chunkRepository.EXPECT().GetAll(context.TODO(), ownerID).Return(emptyChunks, nil)
		lister := chunklister.New(ownerID, factory.GetChunkRepository())
		chunks, err := lister.List(context.TODO())

		assert.Nil(t, err)
		assert.Empty(t, chunks)
	})

	t.Run("propagate the error when the lookup fails", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		ownerID := "owner-id"
		lookupError := errors.New("error listing chunks")
		chunkRepository := factory.GetChunkRepositoryMock()
		chunkRepository.EXPECT().GetAll(context.TODO(), ownerID).Return(nil, lookupError)
		lister := chunklister.New(ownerID, factory.GetChunkRepository())
		chunks, err := lister.List(context.TODO())

		assert.Nil(t, chunks)
		assert.Equal(t, lookupError, err)
	})
}
