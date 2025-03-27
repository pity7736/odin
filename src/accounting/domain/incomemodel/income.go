package incomemodel

import (
	"time"

	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
)

type Income struct {
	id     string
	amount moneymodel.Money
	date   time.Time
}

func New(id string, amount moneymodel.Money, date time.Time) *Income {
	return &Income{id: id, amount: amount, date: date}
}

func (self *Income) ID() string {
	return self.id
}

func (self *Income) Amount() moneymodel.Money {
	return self.amount
}

func (self *Income) Date() time.Time {
	return self.date
}
