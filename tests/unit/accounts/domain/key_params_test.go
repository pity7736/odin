package domain_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"raiseexception.dev/odin/src/accounts/domain/keyparams"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

func TestKeyParamsShould(t *testing.T) {
	t.Run("create valid key params with all fields accessible via getters", func(t *testing.T) {
		params, err := keyparams.New("argon2id", 3, 65536, 4, "random-salt-value")
		assert.Nil(t, err)
		assert.Equal(t, "argon2id", params.Algorithm())
		assert.Equal(t, 3, params.Iterations())
		assert.Equal(t, 65536, params.Memory())
		assert.Equal(t, 4, params.Parallelism())
		assert.Equal(t, "random-salt-value", params.Salt())
	})
	t.Run("reject empty algorithm", func(t *testing.T) {
		params, err := keyparams.New("", 3, 65536, 4, "salt")
		assert.Equal(t, keyparams.KeyParams{}, params)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, "El algoritmo es obligatorio", odinError.ExternalError())
	})
	t.Run("reject non-positive iterations", func(t *testing.T) {
		params, err := keyparams.New("argon2id", 0, 65536, 4, "salt")
		assert.Equal(t, keyparams.KeyParams{}, params)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, "Las iteraciones deben ser positivas", odinError.ExternalError())
	})
	t.Run("reject non-positive memory", func(t *testing.T) {
		params, err := keyparams.New("argon2id", 3, 0, 4, "salt")
		assert.Equal(t, keyparams.KeyParams{}, params)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, "La memoria debe ser positiva", odinError.ExternalError())
	})
	t.Run("reject non-positive parallelism", func(t *testing.T) {
		params, err := keyparams.New("argon2id", 3, 65536, 0, "salt")
		assert.Equal(t, keyparams.KeyParams{}, params)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, "El paralelismo debe ser positivo", odinError.ExternalError())
	})
	t.Run("reject empty salt", func(t *testing.T) {
		params, err := keyparams.New("argon2id", 3, 65536, 4, "")
		assert.Equal(t, keyparams.KeyParams{}, params)
		var odinError *odinerrors.Error
		assert.True(t, errors.As(err, &odinError))
		assert.Equal(t, "La sal es obligatoria", odinError.ExternalError())
	})
}
