package use_cases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/vault/application/use_cases/chunkcreator"
	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
	"raiseexception.dev/odin/tests/unit/testrepositoryfactory"
)

func TestCreateChunkShould(t *testing.T) {
	t.Run("store the chunk when the owner has no chunk with that id", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		id := uuid.NewString()
		ownerID := "owner-id"
		content := "encrypted-content"
		chunkRepository := factory.GetChunkRepositoryMock()
		chunkRepository.EXPECT().GetByID(context.TODO(), ownerID, id).Return(nil, nil)
		var storedChunk *chunkmodel.EncryptedChunk
		chunkRepository.EXPECT().Add(context.TODO(), mock.Anything).Run(func(ctx context.Context, chunk *chunkmodel.EncryptedChunk) {
			storedChunk = chunk
		}).Return(nil)
		creator := chunkcreator.New(id, ownerID, content, factory.GetChunkRepository())
		chunk, err := creator.Create(context.TODO())

		assert.Nil(t, err)
		assert.NotNil(t, chunk)
		assert.Equal(t, id, chunk.ID())
		assert.Equal(t, ownerID, chunk.OwnerID())
		assert.Equal(t, content, chunk.Content())
		assert.Equal(t, chunk, storedChunk)
		chunkRepository.AssertCalled(t, "Add", context.TODO(), mock.Anything)
	})

	t.Run("reject the chunk when the owner already has one with that id", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		id := uuid.NewString()
		ownerID := "owner-id"
		existingChunk, _ := chunkmodel.New(id, ownerID, "existing-content")
		chunkRepository := factory.GetChunkRepositoryMock()
		chunkRepository.EXPECT().GetByID(context.TODO(), ownerID, id).Return(existingChunk, nil)
		creator := chunkcreator.New(id, ownerID, "new-content", factory.GetChunkRepository())
		chunk, err := creator.Create(context.TODO())

		assert.Nil(t, chunk)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, odinerrors.AlreadyExists, odinError.Tag())
		assert.Equal(t, "El elemento ya existe", odinError.ExternalError())
		chunkRepository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("propagate the error when looking up the chunk fails", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		id := uuid.NewString()
		ownerID := "owner-id"
		lookupError := errors.New("error getting chunk")
		chunkRepository := factory.GetChunkRepositoryMock()
		chunkRepository.EXPECT().GetByID(context.TODO(), ownerID, id).Return(nil, lookupError)
		creator := chunkcreator.New(id, ownerID, "content", factory.GetChunkRepository())
		chunk, err := creator.Create(context.TODO())

		assert.Nil(t, chunk)
		assert.Equal(t, lookupError, err)
		chunkRepository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("propagate the error when persisting the chunk fails", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		id := uuid.NewString()
		ownerID := "owner-id"
		persistError := errors.New("error saving chunk")
		chunkRepository := factory.GetChunkRepositoryMock()
		chunkRepository.EXPECT().GetByID(context.TODO(), ownerID, id).Return(nil, nil)
		chunkRepository.EXPECT().Add(context.TODO(), mock.Anything).Return(persistError)
		creator := chunkcreator.New(id, ownerID, "content", factory.GetChunkRepository())
		chunk, err := creator.Create(context.TODO())

		assert.Nil(t, chunk)
		assert.Equal(t, persistError, err)
	})

	t.Run("reject the chunk when the id is not a valid uuid", func(t *testing.T) {
		factory := testrepositoryfactory.New(t)
		ownerID := "owner-id"
		chunkRepository := factory.GetChunkRepositoryMock()
		chunkRepository.EXPECT().GetByID(context.TODO(), ownerID, "not-a-uuid").Return(nil, nil)
		creator := chunkcreator.New("not-a-uuid", ownerID, "content", factory.GetChunkRepository())
		chunk, err := creator.Create(context.TODO())

		assert.Nil(t, chunk)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
		chunkRepository.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})
}
