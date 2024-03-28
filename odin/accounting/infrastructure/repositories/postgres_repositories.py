from typing import Optional

import asyncpg

from odin import settings
from odin.accounting.application.repositories import CategoryRepository, WalletRepository, TransferRepository
from odin.accounting.domain import CategoryType
from odin.accounting.domain.models import Category, Wallet, Income, Expense, Transfer
from odin.accounts.domain import User
from odin.shared.db_connection_manager import DBConnectionManager


class PostgresCategoryRepository(CategoryRepository):

    def __init__(self):
        self._connection_manager = DBConnectionManager()

    async def get_by_name_and_user(self, name: str, user: User) -> Optional[Category]:
        pass

    async def add(self, category: Category):
        async with self._connection_manager as connection:
            await connection.execute(
                'insert into categories (id, name, type, user_id) values ($1, $2, $3, $4)',
                category.id,
                category.name,
                category.type.value,
                category.user.id
            )

    async def get_all_by_user_and_type(self, user: User, type: CategoryType) -> tuple[Category]:
        connection = await asyncpg.connect(
            host=settings.DB_HOST,
            user=settings.DB_USER,
            database=settings.DB_NAME,
            password=settings.DB_PASSWORD,
            port=settings.DB_PORT,
        )
        records = await connection.fetch(
            'select id, name, type from categories where user_id = $1 and type = $2',
            user.id,
            type.value
        )
        await connection.close()
        categories = []
        for record in records:
            categories.append(Category(**record))
        return tuple(categories)

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
