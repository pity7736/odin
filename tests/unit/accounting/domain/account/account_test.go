package account

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
)

func Test_givenNewAccount_WhenIDIsEmpty_ThenReturnError(t *testing.T) {
	balance, _ := moneymodel.New("100")
	_, err := accountmodel.New(
		"",
		"savings",
		"user id",
		balance,
		balance,
	)

	assert.Equal(t, errors.New("id cannot be empty"), err)
}

func Test_givenNewAccount_WhenNameIsEmpty_ThenReturnError(t *testing.T) {
	balance, _ := moneymodel.New("100")
	_, err := accountmodel.New(
		"some id",
		"",
		"user id",
		balance,
		balance,
	)

	assert.Equal(t, errors.New("name cannot be empty"), err)
}

func Test_givenNewAccount_WhenUserIDIsEmpty_ThenReturnError(t *testing.T) {
	balance, _ := moneymodel.New("100")
	_, err := accountmodel.New(
		"some id",
		"savings",
		"",
		balance,
		balance,
	)

	assert.Equal(t, errors.New("user id cannot be empty"), err)
}

func Test_givenNewAccount_WhenNegativeInitialBalance_ThenReturnError(t *testing.T) {
	initialBalance, _ := moneymodel.New("-100")
	balance, _ := moneymodel.New("100")
	_, err := accountmodel.New(
		"some id",
		"savings",
		"user id",
		initialBalance,
		balance,
	)

	assert.Equal(t, errors.New("initial balance must be positive"), err)
}

func Test_givenNewAccount_WhenNegativeBalance_ThenReturnError(t *testing.T) {
	initialBalance, _ := moneymodel.New("100")
	balance, _ := moneymodel.New("-100")
	_, err := accountmodel.New(
		"some id",
		"savings",
		"user id",
		initialBalance,
		balance,
	)

	assert.Equal(t, errors.New("balance must be positive"), err)
}
