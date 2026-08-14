package moneymodel

import (
	"strings"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type Currency struct {
	code string
}

func COP() Currency {
	return Currency{code: "COP"}
}

func USD() Currency {
	return Currency{code: "USD"}
}

func CurrencyFromString(code string) (Currency, error) {
	normalized := strings.ToUpper(code)
	switch normalized {
	case "COP":
		return COP(), nil
	case "USD":
		return USD(), nil
	default:
		return Currency{}, odinerrors.NewErrorBuilder("invalid currency").
			WithExternalMessage("Moneda inválida").
			WithTag(odinerrors.Domain).
			Build()
	}
}

func (self Currency) Equals(other Currency) bool {
	return self == other
}

func (self Currency) String() string {
	return self.code
}
