package money

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

func TestMoneyShould(t *testing.T) {
	t.Run("default to COP when no currency provided", func(t *testing.T) {
		money, err := moneymodel.New("100")
		assert.Nil(t, err)
		assert.True(t, money.Currency().Equals(moneymodel.COP()))
	})

	t.Run("use provided currency", func(t *testing.T) {
		money, err := moneymodel.New("100", moneymodel.USD())
		assert.Nil(t, err)
		assert.True(t, money.Currency().Equals(moneymodel.USD()))
	})

	t.Run("return Spanish domain error for invalid value", func(t *testing.T) {
		_, err := moneymodel.New("not-a-number")
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "El valor ingresado no es un monto válido", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("expose currency via getter", func(t *testing.T) {
		money, _ := moneymodel.New("500", moneymodel.USD())
		assert.True(t, money.Currency().Equals(moneymodel.USD()))
	})
}
