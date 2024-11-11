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
	return &factory{t: t}
}

func (f *factory) GetCategoryRepository() repositories.CategoryRepository {
	return f.GetCategoryRepositoryMock()
}

func (f *factory) GetCategoryRepositoryMock() *mocks.MockCategoryRepository {
	if f.categoryRepository == nil {
		f.categoryRepository = mocks.NewMockCategoryRepository(f.t)
	}
	return f.categoryRepository
}
