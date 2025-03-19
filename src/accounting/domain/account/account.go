package accountmodel

import (
	"errors"
	"time"

	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
)

type Account struct {
	name           string
	initialBalance moneymodel.Money
	userID         string
	id             string
	balance        moneymodel.Money
	createdAt      time.Time
}

func New(id, name, userID string, initialBalance, balance moneymodel.Money, createdAt time.Time) (*Account, error) {
	err := validateData(id, name, userID, initialBalance, balance)
	if err != nil {
		return nil, err
	}
	return &Account{
		id:             id,
		name:           name,
		initialBalance: initialBalance,
		userID:         userID,
		balance:        balance,
		createdAt:      createdAt,
	}, nil
}

func validateData(id, name, userID string, initialBalance, balance moneymodel.Money) error {
	if initialBalance.IsNegative() {
		return errors.New("initial balance must be positive")
	}
	if balance.IsNegative() {
		return errors.New("balance must be positive")
	}
	if id == "" {
		return errors.New("id cannot be empty")
	}
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if userID == "" {
		return errors.New("user id cannot be empty")
	}
	return nil
}

func (self *Account) ID() string {
	return self.id
}

func (self *Account) Name() string {
	return self.name
}

func (self *Account) InitialBalance() moneymodel.Money {
	return self.initialBalance
}

func (self *Account) UserID() string {
	return self.userID
}

func (self *Account) Balance() moneymodel.Money {
	return self.balance
}

func (self *Account) CreatedAt() time.Time {
	return self.createdAt
}
