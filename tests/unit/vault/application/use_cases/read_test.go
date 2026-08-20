package use_cases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/vault/application/use_cases/chunkgetter"
	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestReadChunkShould(t *testing.T) {
	t.Run("return the chunk when the owner has one with that id", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		id := uuid.NewString()
		ownerID := "owner-id"
		storedChunk, _ := chunkmodel.New(id, ownerID, "encrypted-content")
		chunkRepository := factory.GetChunkRepositoryMock()
		chunkRepository.EXPECT().Get(context.TODO(), ownerID, id).Return(storedChunk, nil)
		getter := chunkgetter.New(id, ownerID, factory.GetChunkRepository())
		chunk, err := getter.Get(context.TODO())

		assert.Nil(t, err)
		assert.Equal(t, storedChunk, chunk)
	})

	t.Run("propagate the not found error verbatim when the owner has no chunk with that id", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		id := uuid.NewString()
		ownerID := "owner-id"
		notFoundError := odinerrors.NewErrorBuilder("chunk not found").
			WithExternalMessage("El elemento no existe").
			WithTag(odinerrors.NotFound).
			Build()
		chunkRepository := factory.GetChunkRepositoryMock()
		chunkRepository.EXPECT().Get(context.TODO(), ownerID, id).Return(nil, notFoundError)
		getter := chunkgetter.New(id, ownerID, factory.GetChunkRepository())
		chunk, err := getter.Get(context.TODO())

		assert.Nil(t, chunk)
		assert.Equal(t, notFoundError, err)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, odinerrors.NotFound, odinError.Tag())
	})

	t.Run("propagate the error when the lookup fails", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		id := uuid.NewString()
		ownerID := "owner-id"
		lookupError := errors.New("error getting chunk")
		chunkRepository := factory.GetChunkRepositoryMock()
		chunkRepository.EXPECT().Get(context.TODO(), ownerID, id).Return(nil, lookupError)
		getter := chunkgetter.New(id, ownerID, factory.GetChunkRepository())
		chunk, err := getter.Get(context.TODO())

		assert.Nil(t, chunk)
		assert.Equal(t, lookupError, err)
	})
}
