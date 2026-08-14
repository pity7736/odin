package accountmodel

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"raiseexception.dev/odin/src/accounting/domain/accounttypemodel"
	categorymodel "raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/domain/incomemodel"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
)

type Account struct {
	incomes        []*incomemodel.Income
	initialBalance moneymodel.Money
	balance        moneymodel.Money
	createdAt      time.Time
	accountType    accounttypemodel.AccountType
	name           string
	userID         string
	id             string
}

func New(name, userID string, initialBalance moneymodel.Money, accountType accounttypemodel.AccountType) (*Account, error) {
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
		accountType,
		time.Now(),
	)
}

func NewFromRepository(id, name, userID string, initialBalance, balance moneymodel.Money, accountType accounttypemodel.AccountType, createdAt time.Time) (*Account, error) {
	trimmedName := strings.TrimSpace(name)
	err := validateData(id, trimmedName, userID, initialBalance, balance)
	if err != nil {
		return nil, err
	}
	return &Account{
		id:             id,
		name:           trimmedName,
		initialBalance: initialBalance,
		userID:         userID,
		balance:        balance,
		accountType:    accountType,
		createdAt:      createdAt,
	}, nil
}

func validateData(id, name, userID string, initialBalance, balance moneymodel.Money) error {
	if initialBalance.IsNegative() {
		return odinerrors.NewErrorBuilder("initial balance must not be negative").
			WithTag(odinerrors.Domain).
			WithExternalMessage("El saldo inicial no puede ser negativo").
			Build()
	}
	if balance.IsNegative() {
		return odinerrors.NewErrorBuilder("balance must not be negative").
			WithTag(odinerrors.Domain).
			WithExternalMessage("El saldo no puede ser negativo").
			Build()
	}
	if id == "" {
		return odinerrors.NewErrorBuilder("id cannot be empty").
			WithTag(odinerrors.Domain).
			WithExternalMessage("id cannot be empty").
			Build()
	}
	if name == "" {
		return odinerrors.NewErrorBuilder("name cannot be empty").
			WithTag(odinerrors.Domain).
			WithExternalMessage("El nombre es obligatorio").
			Build()
	}
	if len([]rune(name)) > 255 {
		return odinerrors.NewErrorBuilder("name is too long").
			WithTag(odinerrors.Domain).
			WithExternalMessage("El nombre es demasiado largo").
			Build()
	}
	if userID == "" {
		return odinerrors.NewErrorBuilder("user id cannot be empty").
			WithTag(odinerrors.Domain).
			WithExternalMessage("user id cannot be empty").
			Build()
	}
	return nil
}

func (self *Account) ValidateOwnership(requestContext *requestcontext.RequestContext) error {
	if self.UserID() != requestContext.UserID() {
		return odinerrors.NewErrorBuilder("cuenta no pertenece a usuario logueado").
			WithTag(odinerrors.Domain).
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

func (self *Account) Type() accounttypemodel.AccountType {
	return self.accountType
}

func (self *Account) Currency() moneymodel.Currency {
	return self.balance.Currency()
}

func (self *Account) CreatedAt() time.Time {
	return self.createdAt
}

func (self *Account) CreateIncome(amount moneymodel.Money, date time.Time, category categorymodel.Category) (*incomemodel.Income, error) {
	minimalAmount := moneymodel.MustNew("1")
	if amount.Less(minimalAmount) {
		return nil, odinerrors.NewErrorBuilder("amount error").
			WithTag(odinerrors.Domain).
			WithExternalMessage("el monto debe ser mayor o igual a 1").
			Build()
	}
	if self.createdAt.After(date) {
		return nil, odinerrors.NewErrorBuilder("date error").
			WithTag(odinerrors.Domain).
			WithExternalMessage("la fecha del ingreso debe ser posterior a la fecha de creación de la cuenta").
			Build()
	}
	if !category.IsIncome() {
		return nil, odinerrors.NewErrorBuilder("category error").
			WithTag(odinerrors.Domain).
			WithExternalMessage("la categoría no es de ingreso").
			Build()
	}
	incomeID, _ := uuid.NewV7()
	income := incomemodel.New(incomeID.String(), amount, date)
	self.balance = self.balance.Subtract(amount)
	self.incomes = append(self.incomes, income)
	return income, nil
}

func (self *Account) Incomes() []*incomemodel.Income {
	return self.incomes
}
