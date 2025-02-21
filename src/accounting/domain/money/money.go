package moneymodel

import "github.com/govalues/decimal"

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
		return Money{}, err
	}
	return Money{value: val, currency: currency}, nil
}

func (self Money) IsNegative() bool {
	return self.value.IsNeg()
}
