package testrepositoryfactory

import (
	"testing"

	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/tests/unit/mocks"
)

type Factory struct {
	t                  *testing.T
	categoryRepository *mocks.MockCategoryRepository
}

func New(t *testing.T) *Factory {
	return &Factory{t: t}
}

func (f *Factory) GetCategoryRepository() repositories.CategoryRepository {
	return f.GetCategoryRepositoryMock()
}

func (f *Factory) GetCategoryRepositoryMock() *mocks.MockCategoryRepository {
	if f.categoryRepository == nil {
		f.categoryRepository = mocks.NewMockCategoryRepository(f.t)
	}
	return f.categoryRepository
}
