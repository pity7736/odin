package repositories

import (
	"context"

	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"
	"raiseexception.dev/odin/src/accounting/domain/category"
)

type CategoryRepository interface {
	Add(ctx context.Context, category *categorymodel.Category) error
	GetAll(ctx context.Context, userID string) []*categorymodel.Category
	GetByID(ctx context.Context, id string) (*categorymodel.Category, error)
}

type AccountRepository interface {
	Add(ctx context.Context, account *accountmodel.Account) error
	GetAll(ctx context.Context) ([]*accountmodel.Account, error)
	GetByID(ctx context.Context, id string) (*accountmodel.Account, error)
	Save(ctx context.Context, account *accountmodel.Account) error
}
