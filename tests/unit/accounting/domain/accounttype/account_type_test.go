package accounttype

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"raiseexception.dev/odin/src/accounting/domain/accounttypemodel"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

func TestAccountTypeShould(t *testing.T) {
	t.Run("create savings type from string", func(t *testing.T) {
		accountType, err := accounttypemodel.NewFromString("savings")
		assert.Nil(t, err)
		assert.True(t, accountType.Equals(accounttypemodel.Savings()))
		assert.Equal(t, "savings", accountType.String())
	})

	t.Run("create credit_card type from string", func(t *testing.T) {
		accountType, err := accounttypemodel.NewFromString("credit_card")
		assert.Nil(t, err)
		assert.True(t, accountType.Equals(accounttypemodel.CreditCard()))
		assert.Equal(t, "credit_card", accountType.String())
	})

	t.Run("create cash type from string", func(t *testing.T) {
		accountType, err := accounttypemodel.NewFromString("cash")
		assert.Nil(t, err)
		assert.True(t, accountType.Equals(accounttypemodel.Cash()))
		assert.Equal(t, "cash", accountType.String())
	})

	t.Run("accept case-insensitive input", func(t *testing.T) {
		accountType, err := accounttypemodel.NewFromString("SAVINGS")
		assert.Nil(t, err)
		assert.True(t, accountType.Equals(accounttypemodel.Savings()))
	})

	t.Run("return domain error for invalid type", func(t *testing.T) {
		_, err := accounttypemodel.NewFromString("invalid")
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "Tipo de cuenta inválido", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("return domain error for empty string", func(t *testing.T) {
		_, err := accounttypemodel.NewFromString("")
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "Tipo de cuenta inválido", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("savings does not equal credit_card", func(t *testing.T) {
		assert.False(t, accounttypemodel.Savings().Equals(accounttypemodel.CreditCard()))
	})
}
