package repositoryfactory

import (
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/pgrepositories"
)

type RepositoryFactory interface {
	GetCategoryRepository() repositories.CategoryRepository
}

type repositoryFactory struct{}

func New() RepositoryFactory {
	return &repositoryFactory{}
}

func (r *repositoryFactory) GetCategoryRepository() repositories.CategoryRepository {
	return pgrepositories.NewPGCategoryRepository()
}
