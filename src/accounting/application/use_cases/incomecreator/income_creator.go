package incomecreator

import (
	"context"
	"time"

	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	categorymodel "raiseexception.dev/odin/src/accounting/domain/category"
	"raiseexception.dev/odin/src/accounting/domain/incomemodel"
	moneymodel "raiseexception.dev/odin/src/accounting/domain/money"
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/accountingrepositoryfactory"
	"raiseexception.dev/odin/src/shared/domain/requestcontext"
)

type IncomeCreator struct {
	accountRepository  repositories.AccountRepository
	categoryRepository repositories.CategoryRepository
	accountID          string
	categoryID         string
	amount             moneymodel.Money
	date               time.Time
}

func New(factory accountingrepositoryfactory.RepositoryFactory, amount moneymodel.Money, date time.Time, categoryID string, accountID string) *IncomeCreator {
	return &IncomeCreator{
		accountRepository:  factory.GetAccountRepository(),
		categoryRepository: factory.GetCategoryRepository(),
		accountID:          accountID,
		categoryID:         categoryID,
		amount:             amount,
		date:               date,
	}
}

func (self *IncomeCreator) Create(ctx context.Context) (*incomemodel.Income, error) {
	account, err := self.getAccount(ctx)
	if err != nil {
		return nil, err
	}
	category, err := self.getCategory(ctx)
	if err != nil {
		return nil, err
	}
	return self.createIncome(ctx, err, account, category)
}

func (self *IncomeCreator) getAccount(ctx context.Context) (*accountmodel.Account, error) {
	account, err := self.accountRepository.GetByID(ctx, self.accountID)
	if err != nil {
		return nil, err
	}
	requestContext := ctx.Value(requestcontext.Key).(*requestcontext.RequestContext)
	if err = account.ValidateOwnership(requestContext); err != nil {
		return nil, err
	}
	return account, nil
}

func (self *IncomeCreator) getCategory(ctx context.Context) (*categorymodel.Category, error) {
	category, err := self.categoryRepository.GetByID(ctx, self.categoryID)
	if err != nil {
		return nil, err
	}
	requestContext := ctx.Value(requestcontext.Key).(*requestcontext.RequestContext)
	if err = category.ValidateOwnership(requestContext); err != nil {
		return nil, err
	}
	return category, nil
}

func (self *IncomeCreator) createIncome(ctx context.Context, err error, account *accountmodel.Account, category *categorymodel.Category) (*incomemodel.Income, error) {
	income, err := account.CreateIncome(self.amount, self.date, *category)
	if err != nil {
		return nil, err
	}
	self.accountRepository.Save(ctx, account)
	return income, nil
}
