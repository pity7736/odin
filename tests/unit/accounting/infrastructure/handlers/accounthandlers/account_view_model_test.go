package accounthandlers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	"raiseexception.dev/odin/src/accounting/domain/accounttypemodel"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/accounting/infrastructure/api/handlers/accounthandler/accountviewmodel"
)

func TestAccountViewModelShould(t *testing.T) {
	t.Run("map savings to its Spanish label", func(t *testing.T) {
		account := buildViewModelAccount(t, accounttypemodel.Savings(), moneymodel.COP(), time.Now())
		viewModel := accountviewmodel.New(account)
		assert.Equal(t, "Ahorros", viewModel.Type)
	})

	t.Run("map credit card to its Spanish label", func(t *testing.T) {
		account := buildViewModelAccount(t, accounttypemodel.CreditCard(), moneymodel.COP(), time.Now())
		viewModel := accountviewmodel.New(account)
		assert.Equal(t, "Tarjeta de crédito", viewModel.Type)
	})

	t.Run("map cash to its Spanish label", func(t *testing.T) {
		account := buildViewModelAccount(t, accounttypemodel.Cash(), moneymodel.COP(), time.Now())
		viewModel := accountviewmodel.New(account)
		assert.Equal(t, "Efectivo", viewModel.Type)
	})

	t.Run("format the creation date as ISO", func(t *testing.T) {
		account := buildViewModelAccount(t, accounttypemodel.Savings(), moneymodel.COP(), time.Date(2026, time.August, 13, 10, 0, 0, 0, time.UTC))
		viewModel := accountviewmodel.New(account)
		assert.Equal(t, "2026-08-13", viewModel.CreatedAt)
	})

	t.Run("keep the currency as its code", func(t *testing.T) {
		account := buildViewModelAccount(t, accounttypemodel.Savings(), moneymodel.USD(), time.Now())
		viewModel := accountviewmodel.New(account)
		assert.Equal(t, "USD", viewModel.Currency)
	})

	t.Run("carry the account identity and balances", func(t *testing.T) {
		account := buildViewModelAccount(t, accounttypemodel.Savings(), moneymodel.COP(), time.Now())
		viewModel := accountviewmodel.New(account)
		assert.Equal(t, account.ID(), viewModel.ID)
		assert.Equal(t, account.Name(), viewModel.Name)
		assert.Equal(t, account.InitialBalance().String(), viewModel.InitialBalance)
		assert.Equal(t, account.Balance().String(), viewModel.Balance)
	})
}

func buildViewModelAccount(t *testing.T, accountType accounttypemodel.AccountType, currency moneymodel.Currency, createdAt time.Time) *accountmodel.Account {
	t.Helper()
	balance := moneymodel.MustNew("1000", currency)
	account, err := accountmodel.NewFromRepository("id-1", "test", "user-1", balance, balance, accountType, createdAt)
	assert.NoError(t, err)
	return account
}
