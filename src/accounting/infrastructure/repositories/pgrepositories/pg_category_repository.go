package pgrepositories

import (
	"context"

	"raiseexception.dev/odin/src/accounting/domain/category"
)

type PGCategoryRepository struct {
	categories map[string]*categorymodel.Category
}

func NewPGCategoryRepository() *PGCategoryRepository {
	return &PGCategoryRepository{categories: make(map[string]*categorymodel.Category)}
}

func (self *PGCategoryRepository) Add(ctx context.Context, category *categorymodel.Category) error {
	self.categories[category.ID()] = category
	return nil
}

func (self *PGCategoryRepository) GetAll(ctx context.Context, userID string) []*categorymodel.Category {
	result := make([]*categorymodel.Category, 0, len(self.categories))
	for _, category := range self.categories {
		if category.UserID() == userID {
			result = append(result, category)
		}
	}
	return result
}

func (self *PGCategoryRepository) GetByID(ctx context.Context, id string) (*categorymodel.Category, error) {
	//TODO implement me
	panic("implement me")
}
