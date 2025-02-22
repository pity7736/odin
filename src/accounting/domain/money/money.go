package moneymodel

import (
	"fmt"

	"github.com/govalues/decimal"
)

type Money struct {
	value    decimal.Decimal
	currency *Currency
}

func New(value string, currencies ...*Currency) (Money, error) {
	currency := COP()
	if len(currencies) > 0 {
		currency = currencies[0]
	}
	val, err := decimal.Parse(value)
	if err != nil {
		return Money{}, fmt.Errorf(`%s is not valid money value`, value)
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
