package repositories

import (
	"context"

	"raiseexception.dev/odin/src/accounting/domain/category"
)

type CategoryRepository interface {
	Add(ctx context.Context, category *category.Category) error
}
