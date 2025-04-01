package testrepositoryfactory

import (
	"testing"

	"raiseexception.dev/odin/src/accounting/domain/repositories"
	accountsrepositories "raiseexception.dev/odin/src/accounts/domain/repositories"
	"raiseexception.dev/odin/tests/unit/mocks"
)

type Factory struct {
	t                  *testing.T
	categoryRepository *mocks.MockCategoryRepository
	sessionRepository  *mocks.MockSessionRepository
	userRepository     *mocks.MockUserRepository
	accountsRepository *mocks.MockAccountRepository
	incomeRepository   *mocks.MockIncomeRepository
}

func New(t *testing.T) *Factory {
	return &Factory{t: t}
}

func (self *Factory) GetCategoryRepository() repositories.CategoryRepository {
	return self.GetCategoryRepositoryMock()
}

func (self *Factory) GetCategoryRepositoryMock() *mocks.MockCategoryRepository {
	if self.categoryRepository == nil {
		self.categoryRepository = mocks.NewMockCategoryRepository(self.t)
	}
	return self.categoryRepository
}

func (self *Factory) GetSessionRepository() accountsrepositories.SessionRepository {
	return self.GetSessionRepositoryMock()
}

func (self *Factory) GetSessionRepositoryMock() *mocks.MockSessionRepository {
	if self.sessionRepository == nil {
		self.sessionRepository = mocks.NewMockSessionRepository(self.t)
	}
	return self.sessionRepository
}

func (self *Factory) GetUserRepository() accountsrepositories.UserRepository {
	return self.GetUserRepositoryMock()
}

func (self *Factory) GetUserRepositoryMock() *mocks.MockUserRepository {
	if self.userRepository == nil {
		self.userRepository = mocks.NewMockUserRepository(self.t)
	}
	return self.userRepository
}

func (self *Factory) GetAccountRepository() repositories.AccountRepository {
	return self.GetAccountRepositoryMock()
}

func (self *Factory) GetAccountRepositoryMock() *mocks.MockAccountRepository {
	if self.accountsRepository == nil {
		self.accountsRepository = mocks.NewMockAccountRepository(self.t)
	}
	return self.accountsRepository
}

func (self *Factory) GetIncomeRepository() repositories.IncomeRepository {
	return self.GetIncomeRepositoryMock()
}

func (self *Factory) GetIncomeRepositoryMock() *mocks.MockIncomeRepository {
	if self.incomeRepository == nil {
		self.incomeRepository = mocks.NewMockIncomeRepository(self.t)
	}
	return self.incomeRepository
}
