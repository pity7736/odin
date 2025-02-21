package moneymodel

type Currency struct {
	code string
}

func COP() *Currency {
	return newCurrency("COP")
}

func newCurrency(code string) *Currency {
	return &Currency{code: code}
}
