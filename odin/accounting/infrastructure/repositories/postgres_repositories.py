from typing import Optional

from odin.accounting.application.repositories import CategoryRepository, WalletRepository, TransferRepository
from odin.accounting.domain.models import Category, Wallet, Income, Expense, Transfer
from odin.accounts.domain import User


class PostgresCategoryRepository(CategoryRepository):

    def add(self, category: Category):
        pass

    def get_all_by_user(self, user: User) -> tuple[Category]:
        pass

    def get_by_name(self, name: str) -> Optional[Category]:
        pass


class PostgresWalletRepository(WalletRepository):

    def add(self, wallet: Wallet):
        pass

    def add_expense(self, wallet: Wallet, expense: Expense):
        pass

    def add_income(self, wallet: Wallet, income: Income):
        pass

    def get_by_name(self, name: str) -> Optional[Wallet]:
        pass

    def get_by_name_with_expenses(self, name: str) -> Optional[Wallet]:
        pass

    def get_by_name_with_incomes(self, name: str) -> Optional[Wallet]:
        pass


class PostgresTransferRepository(TransferRepository):

    def add(self, transfer: Transfer):
        pass

    def get_by_id(self, id: str) -> tuple[Transfer]:
        pass
