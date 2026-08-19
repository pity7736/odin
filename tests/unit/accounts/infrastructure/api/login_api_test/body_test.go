package login_api_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"raiseexception.dev/odin/src/accounts/infrastructure/api/loginhandler"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

func TestLoginBodyValidation(t *testing.T) {
	t.Run("valid body returns no error", func(t *testing.T) {
		body := loginhandler.LoginBody{Email: "test@example.com", AuthHash: "some-auth-hash"}
		assert.Nil(t, body.Validate())
	})
	t.Run("empty email returns error", func(t *testing.T) {
		body := loginhandler.LoginBody{Email: "", AuthHash: "some-auth-hash"}
		err := body.Validate()
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, "El correo es obligatorio", odinError.ExternalError())
	})
	t.Run("empty auth hash returns error", func(t *testing.T) {
		body := loginhandler.LoginBody{Email: "test@example.com", AuthHash: ""}
		err := body.Validate()
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, "La contraseña es obligatoria", odinError.ExternalError())
	})
}
