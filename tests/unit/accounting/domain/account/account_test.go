package account

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/tests/builders"
	"raiseexception.dev/odin/tests/builders/categorybuilder"
	"raiseexception.dev/odin/tests/testutils"
)

func Test_givenNewAccount_WhenDataIsValid_ThenReturnANewAccount(t *testing.T) {
	balance, _ := moneymodel.New("100")
	name := "test"
	userID := "1234"
	account, err := accountmodel.New(name, userID, balance)

	assert.Nil(t, err)
	assert.True(t, testutils.IsUUIDv7(account.ID()))
	assert.Equal(t, name, account.Name())
	assert.Equal(t, userID, account.UserID())
	assert.Equal(t, balance, account.Balance())
	assert.Equal(t, balance, account.InitialBalance())
	assert.True(t, testutils.IsTimeClose(time.Now(), account.CreatedAt()))
}

func Test_givenNewAccountFromRepository_WhenIDIsEmpty_ThenReturnError(t *testing.T) {
	balance, _ := moneymodel.New("100")
	_, err := accountmodel.NewFromRepository("", "savings", "user id", balance, balance, time.Now())

	var odinError *odinerrors.Error
	ok := errors.As(err, &odinError)
	assert.True(t, ok)
	assert.Equal(t, "id cannot be empty", odinError.ExternalError())
	assert.Equal(t, odinerrors.DOMAIN, odinError.Tag())
}

func Test_givenNewAccountFromRepository_WhenNameIsEmpty_ThenReturnError(t *testing.T) {
	balance, _ := moneymodel.New("100")
	_, err := accountmodel.NewFromRepository("some id", "", "user id", balance, balance, time.Now())

	var odinError *odinerrors.Error
	ok := errors.As(err, &odinError)
	assert.True(t, ok)
	assert.Equal(t, "name cannot be empty", odinError.ExternalError())
	assert.Equal(t, odinerrors.DOMAIN, odinError.Tag())
}

func Test_givenNewAccountFromRepository_WhenUserIDIsEmpty_ThenReturnError(t *testing.T) {
	balance, _ := moneymodel.New("100")
	_, err := accountmodel.NewFromRepository("some id", "savings", "", balance, balance, time.Now())

	var odinError *odinerrors.Error
	ok := errors.As(err, &odinError)
	assert.True(t, ok)
	assert.Equal(t, "user id cannot be empty", odinError.ExternalError())
	assert.Equal(t, odinerrors.DOMAIN, odinError.Tag())
}

func Test_givenNewAccountFromRepository_WhenNegativeInitialBalance_ThenReturnError(t *testing.T) {
	initialBalance, _ := moneymodel.New("-100")
	balance, _ := moneymodel.New("100")
	_, err := accountmodel.NewFromRepository("some id", "savings", "user id", initialBalance, balance, time.Now())

	var odinError *odinerrors.Error
	ok := errors.As(err, &odinError)
	assert.True(t, ok)
	assert.Equal(t, "initial balance must be positive", odinError.ExternalError())
	assert.Equal(t, odinerrors.DOMAIN, odinError.Tag())
}

func Test_givenNewAccountFromRepository_WhenNegativeBalance_ThenReturnError(t *testing.T) {
	initialBalance, _ := moneymodel.New("100")
	balance, _ := moneymodel.New("-100")
	_, err := accountmodel.NewFromRepository("some id", "savings", "user id", initialBalance, balance, time.Now())

	var odinError *odinerrors.Error
	errors.As(err, &odinError)
	assert.Equal(t, "balance must be positive", odinError.ExternalError())
	assert.Equal(t, odinerrors.DOMAIN, odinError.Tag())
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
		assert.Equal(t, odinerrors.DOMAIN, odinError.Tag())
	})

	t.Run("return error when date is before than account creation", func(t *testing.T) {
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
		assert.Equal(t, odinerrors.DOMAIN, odinError.Tag())
	})

	t.Run("return error when date is before than account creation", func(t *testing.T) {
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
		assert.Equal(t, odinerrors.DOMAIN, odinError.Tag())
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
		assert.Equal(t, odinerrors.DOMAIN, odinError.Tag())
	})

	t.Run("return income and add amount to account balance", func(t *testing.T) {
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
