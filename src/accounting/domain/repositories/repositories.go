package repositories

import (
	"context"

	"raiseexception.dev/odin/src/accounting/domain/category"
)

type CategoryRepository interface {
	Add(ctx context.Context, category *categorymodel.Category) error
	GetAll(ctx context.Context, userID string) []*categorymodel.Category
}
