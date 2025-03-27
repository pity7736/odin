package moneymodel

import (
	"fmt"

	"github.com/govalues/decimal"
	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type Money struct {
	value    decimal.Decimal
	currency *Currency
}

func MustNew(value string, currencies ...*Currency) Money {
	money, err := New(value, currencies...)
	if err != nil {
		panic(err)
	}
	return money
}

func New(value string, currencies ...*Currency) (Money, error) {
	currency := COP()
	if len(currencies) > 0 {
		currency = currencies[0]
	}
	val, err := decimal.Parse(value)
	if err != nil {
		message := fmt.Sprintf(`%s is not valid money value`, value)
		return Money{}, odinerrors.NewErrorBuilder(message).
			WithExternalMessage(message).
			WithTag(odinerrors.DOMAIN).
			Build()
	}
	return Money{value: val, currency: currency}, nil
}

func (self Money) IsNegative() bool {
	return self.value.IsNeg()
}

func (self Money) Value() decimal.Decimal {
	return self.value
}

func (self Money) String() string {
	return self.value.String()
}

func (self Money) Less(amount Money) bool {
	return self.value.Less(amount.value)
}

func (self Money) Subtract(amount Money) Money {
	value, _ := self.value.Sub(amount.value)
	return Money{value: value, currency: self.currency}
}
