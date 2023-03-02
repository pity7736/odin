from typing import Optional

from odin.accounting.models import Wallet, Expense

from .db_client import DBClient


class EdgeDBWalletRepository:

    def __init__(self):
        self._client = DBClient()

    def add(self, wallet: Wallet):
        self._client.execute(
            'insert Wallet {name := <str>$name, balance := <decimal>$balance}',
            name=wallet.name,
            balance=wallet.balance
        )

    def add_expense(self, wallet: Wallet, expense: Expense):
        category_query = 'select Category filter .name = <str>$category_name'
        wallet_query = 'select Wallet filter .name = <str>$wallet_name'
        expense_query = (
            f'insert Expense {{'
            f'date := <cal::local_date>$date, amount := <decimal>$amount, '
            f'category := ({category_query}), wallet := ({wallet_query})}}'
        )
        result = self._client.query_single(
            expense_query,
            category_name=expense.category.name,
            wallet_name=wallet.name,
            date=expense.date,
            amount=expense.amount
        )
        expense.uuid = result.id

    def get_by_name(self, name: str) -> Optional[Wallet]:
        record = self._client.query_single('select Wallet {id, name, balance} filter .name = <str>$name', name=name)
        if record:
            return Wallet(
                name=record.name,
                balance=record.balance,
                uuid=record.id
            )

    def get_by_name_with_expenses(self, name: str) -> Optional[Wallet]:
        expenses_query = 'expenses := .<wallet[is Expense] {id, date, amount}'
        record = self._client.query_single(
            f'select Wallet {{id, name, balance, {expenses_query}}} filter .name = <str>$name',
            name=name
        )
        expenses = [
            Expense(date=expense_data.date, uuid=expense_data.id, amount=expense_data.amount)
            for expense_data in record.expenses
        ]
        if record:
            return Wallet(
                name=record.name,
                balance=record.balance,
                uuid=record.id,
                expenses=expenses
            )
