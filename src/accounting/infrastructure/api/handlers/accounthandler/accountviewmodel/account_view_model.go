package accountviewmodel

import (
	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	"raiseexception.dev/odin/src/accounting/domain/accounttypemodel"
)

type AccountViewModel struct {
	ID             string
	Name           string
	Type           string
	Currency       string
	InitialBalance string
	Balance        string
	CreatedAt      string
}

func New(account *accountmodel.Account) AccountViewModel {
	return AccountViewModel{
		ID:             account.ID(),
		Name:           account.Name(),
		Type:           typeLabel(account.Type()),
		Currency:       account.Currency().String(),
		InitialBalance: account.InitialBalance().String(),
		Balance:        account.Balance().String(),
		CreatedAt:      account.CreatedAt().Format("2006-01-02"),
	}
}

func typeLabel(accountType accounttypemodel.AccountType) string {
	switch accountType {
	case accounttypemodel.Savings():
		return "Ahorros"
	case accounttypemodel.CreditCard():
		return "Tarjeta de crédito"
	case accounttypemodel.Cash():
		return "Efectivo"
	default:
		return ""
	}
}
