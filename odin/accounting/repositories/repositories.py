from abc import ABCMeta, abstractmethod
from typing import Optional

from odin.accounting.models import Category, Wallet, Expense, Income


class CategoryRepository(metaclass=ABCMeta):

    @abstractmethod
    def add(self, category: Category):
        pass

    @abstractmethod
    def get_all(self) -> tuple[Category]:
        pass

    @abstractmethod
    def get_by_name(self, name: str) -> Optional[Category]:
        pass


class WalletRepository(metaclass=ABCMeta):

    @abstractmethod
    def add(self, wallet: Wallet):
        pass

    @abstractmethod
    def add_expense(self, wallet: Wallet, expense: Expense):
        pass

    @abstractmethod
    def add_income(self, wallet: Wallet, income: Income):
        pass

    @abstractmethod
    def get_by_name(self, name: str) -> Optional[Wallet]:
        pass

    @abstractmethod
    def get_by_name_with_expenses(self, name: str) -> Optional[Wallet]:
        pass

    @abstractmethod
    def get_by_name_with_incomes(self, name: str) -> Optional[Wallet]:
        pass
