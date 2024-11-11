package pgrepositories

import (
	"context"

	"raiseexception.dev/odin/src/accounting/domain/category"
)

type PGCategoryRepository struct{}

func NewPGCategoryRepository() *PGCategoryRepository {
	return &PGCategoryRepository{}
}

func (pg *PGCategoryRepository) Add(ctx context.Context, category *category.Category) error {
	return nil
}
