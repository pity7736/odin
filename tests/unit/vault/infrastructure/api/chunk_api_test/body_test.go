package chunk_api_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/vault/infrastructure/api/chunkhandler"
)

func validChunkBody() chunkhandler.CreateChunkBody {
	return chunkhandler.CreateChunkBody{
		ID:      uuid.NewString(),
		Content: "encrypted-content",
	}
}

func TestCreateChunkBodyShould(t *testing.T) {
	t.Run("return no error when the body is valid", func(t *testing.T) {
		body := validChunkBody()
		assert.Nil(t, body.Validate())
	})
	t.Run("return an error when the id is empty", func(t *testing.T) {
		body := validChunkBody()
		body.ID = ""
		err := body.Validate()
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
		assert.Equal(t, "Datos de solicitud inválidos", odinError.ExternalError())
	})
	t.Run("return an error when the content is empty", func(t *testing.T) {
		body := validChunkBody()
		body.Content = ""
		err := body.Validate()
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
		assert.Equal(t, "Datos de solicitud inválidos", odinError.ExternalError())
	})
}
