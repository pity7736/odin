package domain_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/vault/domain/chunkmodel"
)

func TestChunkShould(t *testing.T) {
	t.Run("build a chunk and expose its id, owner id and content", func(t *testing.T) {
		id := uuid.NewString()
		chunk, err := chunkmodel.New(id, "owner-id", "encrypted-content")

		assert.Nil(t, err)
		assert.NotNil(t, chunk)
		assert.Equal(t, id, chunk.ID())
		assert.Equal(t, "owner-id", chunk.OwnerID())
		assert.Equal(t, "encrypted-content", chunk.Content())
	})
	t.Run("reject a chunk whose id is not a valid uuid", func(t *testing.T) {
		chunk, err := chunkmodel.New("not-a-uuid", "owner-id", "encrypted-content")

		assert.Nil(t, chunk)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})
	t.Run("reject a chunk with an empty owner id", func(t *testing.T) {
		chunk, err := chunkmodel.New(uuid.NewString(), "", "encrypted-content")

		assert.Nil(t, chunk)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})
	t.Run("reject a chunk with empty content", func(t *testing.T) {
		chunk, err := chunkmodel.New(uuid.NewString(), "owner-id", "")

		assert.Nil(t, chunk)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})
}
