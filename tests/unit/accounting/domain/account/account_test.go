package account

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
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

	assert.Equal(t, errors.New("id cannot be empty"), err)
}

func Test_givenNewAccountFromRepository_WhenNameIsEmpty_ThenReturnError(t *testing.T) {
	balance, _ := moneymodel.New("100")
	_, err := accountmodel.NewFromRepository("some id", "", "user id", balance, balance, time.Now())

	assert.Equal(t, errors.New("name cannot be empty"), err)
}

func Test_givenNewAccountFromRepository_WhenUserIDIsEmpty_ThenReturnError(t *testing.T) {
	balance, _ := moneymodel.New("100")
	_, err := accountmodel.NewFromRepository("some id", "savings", "", balance, balance, time.Now())

	assert.Equal(t, errors.New("user id cannot be empty"), err)
}

func Test_givenNewAccountFromRepository_WhenNegativeInitialBalance_ThenReturnError(t *testing.T) {
	initialBalance, _ := moneymodel.New("-100")
	balance, _ := moneymodel.New("100")
	_, err := accountmodel.NewFromRepository("some id", "savings", "user id", initialBalance, balance, time.Now())

	assert.Equal(t, errors.New("initial balance must be positive"), err)
}

func Test_givenNewAccountFromRepository_WhenNegativeBalance_ThenReturnError(t *testing.T) {
	initialBalance, _ := moneymodel.New("100")
	balance, _ := moneymodel.New("-100")
	_, err := accountmodel.NewFromRepository("some id", "savings", "user id", initialBalance, balance, time.Now())

	assert.Equal(t, errors.New("balance must be positive"), err)
}
