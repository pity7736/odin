package accountingrepositoryfactory

import (
	"raiseexception.dev/odin/src/accounting/domain/repositories"
	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/pgrepositories"
)

type RepositoryFactory interface {
	GetCategoryRepository() repositories.CategoryRepository
	GetAccountRepository() repositories.AccountRepository
}

type repositoryFactory struct {
	categoryRepository repositories.CategoryRepository
	accountRepository  repositories.AccountRepository
}

func New() RepositoryFactory {
	return &repositoryFactory{
		categoryRepository: pgrepositories.NewPGCategoryRepository(),
		accountRepository:  pgrepositories.NewAccountRepository(),
	}
}

func (self *repositoryFactory) GetCategoryRepository() repositories.CategoryRepository {
	return self.categoryRepository
}

func (self *repositoryFactory) GetAccountRepository() repositories.AccountRepository {
	return self.accountRepository
}
