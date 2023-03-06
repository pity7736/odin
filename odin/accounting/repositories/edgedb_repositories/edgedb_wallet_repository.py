from odin.accounting.models import Wallet, Expense, Category, Income

from .db_client import DBClient
from ..repositories import WalletRepository


class EdgeDBWalletRepository(WalletRepository):

    def __init__(self):
        self._client = DBClient()

    def add(self, wallet):
        self._client.execute(
            'insert Wallet {name := <str>$name, balance := <decimal>$balance}',
            name=wallet.name,
            balance=wallet.balance
        )

    def add_expense(self, wallet, expense):
        return self._add_movement(wallet, expense, "expense")

    def add_income(self, wallet, income):
        return self._add_movement(wallet, income, "income")

    def _add_movement(self, wallet, movement, movement_type):
        category_query = 'select Category filter .name = <str>$category_name'
        wallet_query = 'select Wallet filter .name = <str>$wallet_name'
        expense_query = (
            f'insert Movement {{'
            f'date := <cal::local_date>$date, amount := <decimal>$amount, type := <str>$movement_type, '
            f'category := ({category_query}), wallet := ({wallet_query})}}'
        )
        result = self._client.query_single(
            expense_query,
            category_name=movement.category.name,
            wallet_name=wallet.name,
            date=movement.date,
            amount=movement.amount,
            movement_type=movement_type
        )
        movement.uuid = result.id
        self._update_wallet_balance(wallet)

    def _update_wallet_balance(self, wallet):
        self._client.execute(
            'update Wallet filter .name = <str>$name set {balance := <decimal>$balance}',
            name=wallet.name,
            balance=wallet.balance
        )

    def get_by_name(self, name):
        record = self._client.query_single('select Wallet {id, name, balance} filter .name = <str>$name', name=name)
        if record:
            return Wallet(
                name=record.name,
                balance=record.balance,
                uuid=record.id
            )

    def get_by_name_with_expenses(self, name: str):
        expenses_query = 'expenses := .<wallet[is Movement] {id, date, amount, type, category: {name}}'
        record = self._client.query_single(
            f'select Wallet {{id, name, balance, {expenses_query}}} filter .name = <str>$name',
            name=name
        )
        expenses = []
        for expense_data in record.expenses:
            if expense_data.type == 'expense':
                expenses.append(Expense(
                    date=expense_data.date,
                    uuid=expense_data.id,
                    amount=expense_data.amount,
                    category=Category(name=expense_data.category.name)
                ))
        if record:
            return Wallet(
                name=record.name,
                balance=record.balance,
                uuid=record.id,
                expenses=expenses
            )

    def get_by_name_with_incomes(self, name):
        incomes_query = 'incomes := .<wallet[is Movement] {id, date, amount, type, category: {name}}'
        record = self._client.query_single(
            f'select Wallet {{id, name, balance, {incomes_query}}} filter .name = <str>$name',
            name=name
        )
        incomes = []
        for income_data in record.incomes:
            if income_data.type == 'income':
                incomes.append(Income(
                    date=income_data.date,
                    uuid=income_data.id,
                    amount=income_data.amount,
                    category=Category(name=income_data.category.name)
                ))
        if record:
            return Wallet(
                name=record.name,
                balance=record.balance,
                uuid=record.id,
                incomes=incomes
            )
