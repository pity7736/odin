package money

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

func TestCurrencyShould(t *testing.T) {
	t.Run("create COP from string", func(t *testing.T) {
		currency, err := moneymodel.CurrencyFromString("COP")
		assert.Nil(t, err)
		assert.True(t, currency.Equals(moneymodel.COP()))
		assert.Equal(t, "COP", currency.String())
	})

	t.Run("create USD from string", func(t *testing.T) {
		currency, err := moneymodel.CurrencyFromString("USD")
		assert.Nil(t, err)
		assert.True(t, currency.Equals(moneymodel.USD()))
		assert.Equal(t, "USD", currency.String())
	})

	t.Run("accept case-insensitive input", func(t *testing.T) {
		currency, err := moneymodel.CurrencyFromString("cop")
		assert.Nil(t, err)
		assert.True(t, currency.Equals(moneymodel.COP()))
	})

	t.Run("return domain error for invalid currency", func(t *testing.T) {
		_, err := moneymodel.CurrencyFromString("EUR")
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "Moneda inválida", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("return domain error for empty string", func(t *testing.T) {
		_, err := moneymodel.CurrencyFromString("")
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "Moneda inválida", odinError.ExternalError())
	})

	t.Run("COP equals COP by value", func(t *testing.T) {
		assert.True(t, moneymodel.COP().Equals(moneymodel.COP()))
	})

	t.Run("COP does not equal USD", func(t *testing.T) {
		assert.False(t, moneymodel.COP().Equals(moneymodel.USD()))
	})
}
