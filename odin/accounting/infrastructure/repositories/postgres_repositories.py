from enum import Enum
from typing import Optional


from odin.accounting.application.repositories import CategoryRepository, WalletRepository, TransferRepository
from odin.accounting.domain import CategoryType
from odin.accounting.domain.models import Category, Wallet, Income, Expense, Transfer
from odin.accounts.domain import User
from odin.shared.db_connection_manager import DBConnectionManager


class PostgresCategoryRepository(CategoryRepository):

    def __init__(self):
        self._connection_manager = DBConnectionManager()

    async def get_by_name_and_user(self, name: str, user: User) -> Optional[Category]:
        async with self._connection_manager as connection:
            record = await connection.fetchrow(
                'select id, name, type from categories where name = $1 and user_id = $2',
                name,
                user.id
            )
        if record:
            return Category(**record)

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
        async with self._connection_manager as connection:
            records = await connection.fetch(
                'select id, name, type from categories where user_id = $1 and type = $2',
                user.id,
                type.value
            )
        categories = []
        for record in records:
            categories.append(Category(**record))
        return tuple(categories)

    def get_by_name(self, name: str) -> Optional[Category]:
        pass


class PostgresWalletRepository(WalletRepository):

    _table_name = 'wallets'

    class _MovementType(Enum):
        EXPENSE = 'E'
        INCOME = 'I'

    def __init__(self):
        self._connection_manager = DBConnectionManager()

    async def add(self, wallet: Wallet):
        async with self._connection_manager as connection:
            await connection.execute(
                f'insert into {self._table_name} (id, balance, name, user_id) values ($1, $2, $3, $4)',
                wallet.id,
                wallet.balance,
                wallet.name,
                wallet.user.id
            )

    async def add_expense(self, wallet: Wallet, expense: Expense):
        return await self._add_movement(wallet, expense, self._MovementType.EXPENSE)

    async def add_income(self, wallet: Wallet, income: Income):
        return await self._add_movement(wallet, income, self._MovementType.INCOME)

    async def _add_movement(self, wallet: Wallet, movement: Expense | Income, movement_type: _MovementType):
        async with self._connection_manager as connection:
            await connection.execute(
                '''
                    insert into movements (id, amount, date, movement_type, wallet_id, category_id)
                    values ($1, $2, $3, $4, $5, $6)
                ''',
                movement.id,
                movement.amount,
                movement.date,
                movement_type.value,
                wallet.id,
                movement.category.id
            )

    async def get_by_name(self, name: str) -> Optional[Wallet]:
        async with self._connection_manager as connection:
            record = await connection.fetchrow(
                f'select id, name, balance from {self._table_name} where name = $1',
                name
            )
        if record:
            return Wallet(**record)

    async def get_by_name_with_expenses(self, name: str) -> Optional[Wallet]:
        async with self._connection_manager as connection:
            records = await connection.fetch(
                '''
                    select
                        e.id as e_id, e.amount, e.date, w.balance, w.id as w_id, e.movement_type
                    from wallets as w
                    left join movements as e on (e.wallet_id = w.id)
                    where w.name = $1
                ''',
                name
            )

        expenses = []
        for expense_record in records:
            if expense_record['movement_type'] == 'E':
                expenses.append(Expense(
                    date=expense_record['date'],
                    amount=expense_record['amount'],
                    id=expense_record['e_id'],
                ))
        if records:
            return Wallet(
                id=records[0]['w_id'],
                balance=records[0]['balance'],
                name=name,
                expenses=expenses
            )

    async def get_by_name_with_incomes(self, name: str) -> Optional[Wallet]:
        async with self._connection_manager as connection:
            records = await connection.fetch(
                '''
                    select
                        e.id as e_id, e.amount, e.date, w.balance, w.id as w_id, e.movement_type
                    from wallets as w
                    left join movements as e on (e.wallet_id = w.id)
                    where w.name = $1
                ''',
                name
            )

        expenses = []
        for expense_record in records:
            if expense_record['movement_type'] == 'I':
                expenses.append(Expense(
                    date=expense_record['date'],
                    amount=expense_record['amount'],
                    id=expense_record['e_id'],
                ))
        if records:
            return Wallet(
                id=records[0]['w_id'],
                balance=records[0]['balance'],
                name=name,
                expenses=expenses
            )


class PostgresTransferRepository(TransferRepository):

    _table_name = 'transfers'

    def __init__(self):
        self._connection_manager = DBConnectionManager()

    async def add(self, transfer: Transfer):
        async with self._connection_manager as connection:
            await connection.execute(
                f'''
                    insert into {self._table_name} (id, amount, date, source_id, target_id, expense_id, income_id)
                    values ($1, $2, $3, $4, $5, $6, $7)
                ''',
                transfer.id,
                transfer.amount,
                transfer.date,
                transfer.source.id,
                transfer.target.id,
                transfer.expense.id,
                transfer.income.id
            )

    async def get_by_id(self, id: str) -> Optional[Transfer]:
        async with self._connection_manager as connection:
            record = await connection.fetchrow(
                f'''
                    select
                        t.id as t_id,
                        t.amount,
                        t.date,
                        t.source_id as s_id,
                        s.name as s_name,
                        s.balance as s_balance,
                        t.target_id as ta_id,
                        ta.name as ta_name,
                        ta.balance as ta_balance
                    from {self._table_name} as t
                    join wallets as s on (t.source_id = s.id)
                    join wallets as ta on (t.target_id = ta.id)
                    where t.id = $1
                ''',
                id
            )
        if record:
            return Transfer(
                id=record['t_id'],
                amount=record['amount'],
                date=record['date'],
                source=Wallet(
                    name=record['s_name'],
                    id=record['s_id'],
                    balance=record['s_balance'],
                ),
                target=Wallet(
                    name=record['ta_name'],
                    id=record['ta_id'],
                    balance=record['ta_balance'],
                ),
            )
