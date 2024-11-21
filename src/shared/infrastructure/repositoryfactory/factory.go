package repositoryfactory

import (
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/pgrepositories"
)

type RepositoryFactory interface {
	GetCategoryRepository() repositories.CategoryRepository
}

type repositoryFactory struct {
	categoryRepository repositories.CategoryRepository
}

func New() RepositoryFactory {
	return &repositoryFactory{categoryRepository: pgrepositories.NewPGCategoryRepository()}
}

func (self *repositoryFactory) GetCategoryRepository() repositories.CategoryRepository {
	return self.categoryRepository
}
