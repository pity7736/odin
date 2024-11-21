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

func (self *PGCategoryRepository) Add(ctx context.Context, category *category.Category) error {
	self.categories[category.ID()] = category
	return nil
}

func (self *PGCategoryRepository) GetAll(ctx context.Context) []*category.Category {
	result := make([]*category.Category, 0, len(self.categories))
	for _, category := range self.categories {
		result = append(result, category)
	}
	return result
}
