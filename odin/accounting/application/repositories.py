from abc import ABCMeta, abstractmethod
from typing import Optional

from odin.accounting.domain import CategoryType
from odin.accounting.domain.models import Category, Wallet, Expense, Income, Transfer
from odin.accounts.domain import User


class CategoryRepository(metaclass=ABCMeta):

    @abstractmethod
    async def add(self, category: Category):
        pass

    @abstractmethod
    async def get_all_by_user_and_type(self, user: User, type: CategoryType) -> tuple[Category]:
        pass

    @abstractmethod
    async def get_by_name(self, name: str) -> Optional[Category]:
        pass

    @abstractmethod
    async def get_by_name_and_user(self, name: str, user: User) -> Optional[Category]:
        pass


class WalletRepository(metaclass=ABCMeta):

    @abstractmethod
    async def add(self, wallet: Wallet):
        pass

    @abstractmethod
    async def add_expense(self, wallet: Wallet, expense: Expense):
        pass

    @abstractmethod
    async def add_income(self, wallet: Wallet, income: Income):
        pass

    @abstractmethod
    async def get_by_name(self, name: str) -> Optional[Wallet]:
        pass

    @abstractmethod
    async def get_by_name_with_expenses(self, name: str) -> Optional[Wallet]:
        pass

    @abstractmethod
    async def get_by_name_with_incomes(self, name: str) -> Optional[Wallet]:
        pass


class TransferRepository(metaclass=ABCMeta):

    @abstractmethod
    async def add(self, transfer: Transfer):
        pass

    @abstractmethod
    async def get_by_id(self, id: str) -> Optional[Transfer]:
        pass
