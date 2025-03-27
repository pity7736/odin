package accountmodel

import (
	"time"

	"github.com/google/uuid"
	categorymodel "raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/domain/incomemodel"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
)

type Account struct {
	name           string
	initialBalance moneymodel.Money
	userID         string
	id             string
	balance        moneymodel.Money
	createdAt      time.Time
	incomes        []*incomemodel.Income
}

func New(name, userID string, initialBalance moneymodel.Money) (*Account, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return NewFromRepository(
		id.String(),
		name,
		userID,
		initialBalance,
		initialBalance,
		time.Now(),
	)
}

func NewFromRepository(id, name, userID string, initialBalance, balance moneymodel.Money, createdAt time.Time) (*Account, error) {
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
		return odinerrors.NewErrorBuilder("initial balance must be positive").
			WithTag(odinerrors.DOMAIN).
			WithExternalMessage("initial balance must be positive").
			Build()
	}
	if balance.IsNegative() {
		return odinerrors.NewErrorBuilder("balance must be positive").
			WithTag(odinerrors.DOMAIN).
			WithExternalMessage("balance must be positive").
			Build()
	}
	if id == "" {
		return odinerrors.NewErrorBuilder("id cannot be empty").
			WithTag(odinerrors.DOMAIN).
			WithExternalMessage("id cannot be empty").
			Build()
	}
	if name == "" {
		return odinerrors.NewErrorBuilder("name cannot be empty").
			WithTag(odinerrors.DOMAIN).
			WithExternalMessage("name cannot be empty").
			Build()
	}
	if userID == "" {
		return odinerrors.NewErrorBuilder("user id cannot be empty").
			WithTag(odinerrors.DOMAIN).
			WithExternalMessage("user id cannot be empty").
			Build()
	}
	return nil
}

func (self *Account) ValidateOwnership(requestContext *requestcontext.RequestContext) error {
	if self.UserID() != requestContext.UserID() {
		return odinerrors.NewErrorBuilder("cuenta no pertenece a usuario logueado").
			WithTag(odinerrors.DOMAIN).
			WithExternalMessage("la cuenta no pertenece al usuario logueado").
			Build()
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

func (self *Account) CreateIncome(amount moneymodel.Money, date time.Time, category categorymodel.Category) (*incomemodel.Income, error) {
	minimalAmount := moneymodel.MustNew("1")
	if amount.Less(minimalAmount) {
		return nil, odinerrors.NewErrorBuilder("amount error").
			WithTag(odinerrors.DOMAIN).
			WithExternalMessage("el monto debe ser mayor o igual a 1").
			Build()
	}
	if self.createdAt.After(date) {
		return nil, odinerrors.NewErrorBuilder("date error").
			WithTag(odinerrors.DOMAIN).
			WithExternalMessage("la fecha del ingreso debe ser posterior a la fecha de creación de la cuenta").
			Build()
	}
	if !category.IsIncome() {
		return nil, odinerrors.NewErrorBuilder("category error").
			WithTag(odinerrors.DOMAIN).
			WithExternalMessage("la categoría no es de ingreso").
			Build()
	}
	// TODO: handle error
	incomeID, _ := uuid.NewV7()
	income := incomemodel.New(incomeID.String(), amount, date)
	self.balance = self.balance.Subtract(amount)
	self.incomes = append(self.incomes, income)
	return income, nil
}

func (self *Account) Incomes() []*incomemodel.Income {
	return self.incomes
}
