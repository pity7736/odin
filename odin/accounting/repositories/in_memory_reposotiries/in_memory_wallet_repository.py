import uuid
from typing import Optional

from odin.accounting.models import Wallet
from .in_memory_expense_repository import InMemoryExpenseRepository


class InMemoryWalletRepository:

    _wallets: dict[str, dict] = {}

    def add_expense(self, wallet, expense):
        wallet = self.get_by_name_with_expenses(wallet.name)
        wallet.add_expense(expense)
        expense.uuid = uuid.uuid4()
        self.add(wallet)
        repository = InMemoryExpenseRepository()
        repository.add(expense)

    def add(self, wallet: Wallet):
        self.__class__._wallets[wallet.name] = {
            'name': wallet.name,
            'balance': wallet.balance,
            'expenses_uuid': [expense.uuid for expense in wallet.expenses]
        }

    def update(self, wallet: Wallet):
        if self.get_by_name(wallet.name):
            self.add(wallet)

    def get_by_name(self, name: str) -> Wallet:
        wallet_data = self._wallets.get(name)
        if wallet_data:
            return Wallet(**wallet_data)

    def get_by_name_with_expenses(self, name) -> Optional[Wallet]:
        wallet_data = self._wallets.get(name)
        expenses = []
        repository = InMemoryExpenseRepository()
        for expense_uuid in wallet_data['expenses_uuid']:
            expense = repository.get_by(uuid=expense_uuid)
            expenses.append(expense)
        return Wallet(**wallet_data, expenses=expenses)
