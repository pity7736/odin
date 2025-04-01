package pgrepositories

import (
	"context"

	"raiseexception.dev/odin/src/accounting/domain/incomemodel"
)

type PGIncomeRepository struct {
}

func NewPGIncomeRepository() *PGIncomeRepository {
	return &PGIncomeRepository{}
}

func (self *PGIncomeRepository) Add(ctx context.Context, income *incomemodel.Income) error {
	//TODO implement me
	panic("implement me")
}
