package account

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	"raiseexception.dev/odin/src/accounting/domain/accounttypemodel"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/categorybuilder"
	"raiseexception.dev/odin/tests/testutils"
)

func TestNewAccountShould(t *testing.T) {
	t.Run("return a valid account when data is correct", func(t *testing.T) {
		balance, _ := moneymodel.New("100")
		name := "test"
		userID := "1234"
		account, err := accountmodel.New(name, userID, balance, accounttypemodel.Savings())
		assert.Nil(t, err)
		assert.True(t, testutils.IsUUIDv7(account.ID()))
		assert.Equal(t, name, account.Name())
		assert.Equal(t, userID, account.UserID())
		assert.Equal(t, balance, account.Balance())
		assert.Equal(t, balance, account.InitialBalance())
		assert.True(t, testutils.IsTimeClose(time.Now(), account.CreatedAt()))
		assert.True(t, account.Type().Equals(accounttypemodel.Savings()))
		assert.True(t, account.Currency().Equals(moneymodel.COP()))
	})

	t.Run("trim leading and trailing spaces from name", func(t *testing.T) {
		balance, _ := moneymodel.New("100")
		account, err := accountmodel.New("  my account  ", "1234", balance, accounttypemodel.Savings())
		assert.Nil(t, err)
		assert.Equal(t, "my account", account.Name())
	})

	t.Run("return error when name is only spaces", func(t *testing.T) {
		balance, _ := moneymodel.New("100")
		_, err := accountmodel.New("   ", "1234", balance, accounttypemodel.Savings())
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "El nombre es obligatorio", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("return error when name is too long", func(t *testing.T) {
		balance, _ := moneymodel.New("100")
		longName := string(make([]rune, 256))
		_, err := accountmodel.New(longName, "1234", balance, accounttypemodel.Savings())
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "El nombre es demasiado largo", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})
}

func TestNewAccountFromRepositoryShould(t *testing.T) {
	t.Run("return error when id is empty", func(t *testing.T) {
		balance, _ := moneymodel.New("100")
		_, err := accountmodel.NewFromRepository("", "savings", "user id", balance, balance, accounttypemodel.Savings(), time.Now())
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "id cannot be empty", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("return error when name is empty", func(t *testing.T) {
		balance, _ := moneymodel.New("100")
		_, err := accountmodel.NewFromRepository("some id", "", "user id", balance, balance, accounttypemodel.Savings(), time.Now())
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "El nombre es obligatorio", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("return error when user id is empty", func(t *testing.T) {
		balance, _ := moneymodel.New("100")
		_, err := accountmodel.NewFromRepository("some id", "savings", "", balance, balance, accounttypemodel.Savings(), time.Now())
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "user id cannot be empty", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("return error when initial balance is negative", func(t *testing.T) {
		initialBalance, _ := moneymodel.New("-100")
		balance, _ := moneymodel.New("100")
		_, err := accountmodel.NewFromRepository("some id", "savings", "user id", initialBalance, balance, accounttypemodel.Savings(), time.Now())
		var odinError *odinerrors.Error
		ok := errors.As(err, &odinError)
		assert.True(t, ok)
		assert.Equal(t, "El saldo inicial no puede ser negativo", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("return error when balance is negative", func(t *testing.T) {
		initialBalance, _ := moneymodel.New("100")
		balance, _ := moneymodel.New("-100")
		_, err := accountmodel.NewFromRepository("some id", "savings", "user id", initialBalance, balance, accounttypemodel.Savings(), time.Now())
		var odinError *odinerrors.Error
		errors.As(err, &odinError)
		assert.Equal(t, "El saldo no puede ser negativo", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})
}

func TestCreateIncomeShould(t *testing.T) {
	t.Run("return error when amount is less than 1", func(t *testing.T) {
		account := builders.NewAccountBuilder().Build()
		amount, _ := moneymodel.New("0")
		category := categorybuilder.New().Build()
		income, err := account.CreateIncome(amount, time.Now(), *category)
		var odinError *odinerrors.Error
		errors.As(err, &odinError)
		assert.Nil(t, income)
		assert.Equal(t, "el monto debe ser mayor o igual a 1", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("return error when date is before account creation", func(t *testing.T) {
		account := builders.NewAccountBuilder().Build()
		amount, _ := moneymodel.New("1000")
		category := categorybuilder.New().Build()
		income, err := account.CreateIncome(
			amount,
			account.CreatedAt().AddDate(0, 0, -1),
			*category,
		)
		var odinError *odinerrors.Error
		errors.As(err, &odinError)
		assert.Nil(t, income)
		assert.Equal(t, "la fecha del ingreso debe ser posterior a la fecha de creación de la cuenta", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("return error when category is not income type", func(t *testing.T) {
		account := builders.NewAccountBuilder().Build()
		amount, _ := moneymodel.New("1000")
		category := categorybuilder.New().Build()
		income, err := account.CreateIncome(
			amount,
			account.CreatedAt(),
			*category,
		)
		var odinError *odinerrors.Error
		errors.As(err, &odinError)
		assert.Nil(t, income)
		assert.Equal(t, "la categoría no es de ingreso", odinError.ExternalError())
		assert.Equal(t, odinerrors.Domain, odinError.Tag())
	})

	t.Run("add income and update balance", func(t *testing.T) {
		account := builders.NewAccountBuilder().Build()
		initialBalance := account.InitialBalance()
		amount, _ := moneymodel.New("1000")
		incomeDate := time.Now()
		category := categorybuilder.New().WithIncomeType().Build()
		income, err := account.CreateIncome(amount, incomeDate, *category)
		assert.Nil(t, err)
		assert.True(t, testutils.IsUUIDv7(income.ID()))
		assert.Equal(t, amount, income.Amount())
		assert.Equal(t, incomeDate, income.Date())
		assert.Equal(t, initialBalance.Subtract(amount), account.Balance())
		assert.Equal(t, initialBalance, account.InitialBalance())
		assert.Contains(t, account.Incomes(), income)
	})
}
