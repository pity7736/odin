package register_api_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"raiseexception.dev/odin/src/accounts/infrastructure/api/registerhandler"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

func validBody() registerhandler.RegisterBody {
	return registerhandler.RegisterBody{
		Email:              "test@example.com",
		AuthHash:           "some-auth-hash",
		EncryptedMasterKey: "encrypted-master-key-base64",
	}
}

func TestRegisterBodyShould(t *testing.T) {
	t.Run("valid body returns no error", func(t *testing.T) {
		body := validBody()
		assert.Nil(t, body.Validate())
	})
	t.Run("empty email returns error", func(t *testing.T) {
		body := validBody()
		body.Email = ""
		err := body.Validate()
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, "El correo es obligatorio", odinError.ExternalError())
	})
	t.Run("empty auth hash returns error", func(t *testing.T) {
		body := validBody()
		body.AuthHash = ""
		err := body.Validate()
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, "La contraseña es obligatoria", odinError.ExternalError())
	})
	t.Run("empty encrypted master key returns error", func(t *testing.T) {
		body := validBody()
		body.EncryptedMasterKey = ""
		err := body.Validate()
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, "encrypted master key is required", odinError.Error()[:len("encrypted master key is required")])
		assert.Equal(t, "Datos de solicitud inválidos", odinError.ExternalError())
	})
}
