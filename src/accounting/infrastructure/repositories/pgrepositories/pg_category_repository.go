package pgrepositories

import (
	"context"
	"fmt"

	"raiseexception.dev/odin/src/accounting/domain/category"
)

type PGCategoryRepository struct {
	categories map[string]*category.Category
}

func NewPGCategoryRepository() *PGCategoryRepository {
	m := make(map[string]*category.Category)
	return &PGCategoryRepository{categories: m}
}

func (pg *PGCategoryRepository) Add(ctx context.Context, category *category.Category) error {
	pg.categories[category.ID()] = category
	println("adding category", category)
	return nil
}

func (pg *PGCategoryRepository) GetAll(ctx context.Context) []*category.Category {
	println("getall")
	result := make([]*category.Category, 0, len(pg.categories))
	for _, category := range pg.categories {
		fmt.Printf("category %v", category)
		result = append(result, category)
	}
	return result
}
