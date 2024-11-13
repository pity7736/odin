package pgrepositories

import (
	"context"

	"raiseexception.dev/odin/src/accounting/domain/category"
)

type PGCategoryRepository struct {
	categories map[string]*category.Category
}

func NewPGCategoryRepository() *PGCategoryRepository {
	return &PGCategoryRepository{categories: make(map[string]*category.Category)}
}

func (pg *PGCategoryRepository) Add(ctx context.Context, category *category.Category) error {
	pg.categories[category.ID()] = category
	return nil
}

func (pg *PGCategoryRepository) GetAll(ctx context.Context) []*category.Category {
	result := make([]*category.Category, 0, len(pg.categories))
	for _, category := range pg.categories {
		result = append(result, category)
	}
	return result
}
