package accounttypemodel

import (
	"strings"

	"raiseexception.dev/odin/src/shared/domain/odinerrors"
)

type AccountType struct {
	code string
}

func Savings() AccountType {
	return AccountType{code: "savings"}
}

func CreditCard() AccountType {
	return AccountType{code: "credit_card"}
}

func Cash() AccountType {
	return AccountType{code: "cash"}
}

func NewFromString(value string) (AccountType, error) {
	normalized := strings.ToLower(value)
	switch normalized {
	case "savings":
		return Savings(), nil
	case "credit_card":
		return CreditCard(), nil
	case "cash":
		return Cash(), nil
	default:
		return AccountType{}, odinerrors.NewErrorBuilder("invalid account type").
			WithExternalMessage("Tipo de cuenta inválido").
			WithTag(odinerrors.Domain).
			Build()
	}
}

func (self AccountType) String() string {
	switch self.code {
	case "savings":
		return "savings"
	case "credit_card":
		return "credit_card"
	case "cash":
		return "cash"
	default:
		return ""
	}
}

func (self AccountType) Equals(other AccountType) bool {
	return self.code == other.code
}
