package testrepositoryfactory

import (
	"testing"

	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/tests/unit/mocks"
)

type factory struct {
	t                  *testing.T
	categoryRepository *mocks.MockCategoryRepository
}

func New(t *testing.T) *factory {
	categoryRepository := mocks.NewMockCategoryRepository(t)
	return &factory{
		t:                  t,
		categoryRepository: categoryRepository,
	}
}

func (f *factory) GetCategoryRepository() repositories.CategoryRepository {
	return f.categoryRepository
}

func (f *factory) GetCategoryRepositoryMock() *mocks.MockCategoryRepository {
	return f.categoryRepository
}
