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

    async def get_by_id_and_user(self, id: str, user: User) -> Optional[Category]:
        async with self._connection_manager as connection:
            record = await connection.fetchrow(
                'select id, name, type from categories where id = $1 and user_id = $2',
                id,
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
                getattr(category.user, 'id', None)
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

    async def get_by_name(self, name: str) -> Optional[Category]:
        async with self._connection_manager as connection:
            record = await connection.fetchrow(
                'select id, name, type from categories where name = $1 and user_id is null',
                name,
            )
        if record:
            return Category(**record)


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

    async def get_by_id(self, id: str) -> Optional[Wallet]:
        async with self._connection_manager as connection:
            record = await connection.fetchrow(
                f'''
                    select id, name, balance
                    from {self._table_name}
                    where id = $1
                ''',
                id
            )
        if record:
            return Wallet(**record)

    async def get_expenses_by_wallet_id(self, wallet_id: str) -> list[Expense]:
        async with self._connection_manager as connection:
            records = await connection.fetch(
                '''
                    select
                        m.id as m_id, amount, date, c.id as c_id, c.name
                    from movements as m
                    join categories as c on (m.category_id = c.id)
                    where wallet_id = $1 and m.movement_type = 'E'
                ''',
                wallet_id,
            )
            expenses = []
            for record in records:
                expenses.append(Expense(
                    id=record['m_id'],
                    category=Category(id=record['c_id'], name=record['name']),
                    **record
                ))
            return expenses

    async def get_incomes_by_wallet_id(self, wallet_id: str) -> list[Income]:
        async with self._connection_manager as connection:
            records = await connection.fetch(
                '''
                    select
                        m.id as m_id, amount, date, c.id as c_id, c.name
                    from movements as m
                    join categories as c on (m.category_id = c.id)
                    where wallet_id = $1 and m.movement_type = 'I'
                ''',
                wallet_id,
            )
            incomes = []
            for record in records:
                incomes.append(Income(
                    id=record['m_id'],
                    category=Category(id=record['c_id'], name=record['name']),
                    **record
                ))
            return incomes

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

    async def get_expense_by_wallet_and_expense_id(self, wallet_id: str, expense_id) -> Optional[Expense]:
        async with self._connection_manager as connection:
            record = await connection.fetchrow(
                '''
                    select
                        m.id as m_id, amount, date, c.id as c_id, c.name
                    from movements as m
                    join categories as c on (m.category_id = c.id)
                    where wallet_id = $1 and m.id = $2
                ''',
                wallet_id,
                expense_id
            )
            if record:
                return Expense(
                    id=record['m_id'],
                    category=Category(id=record['c_id'], name=record['name']),
                    **record
                )

    async def get_income_by_wallet_and_income_id(self, wallet_id: str, income_id) -> Optional[Income]:
        pass

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
