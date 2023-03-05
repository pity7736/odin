import uuid

from odin.accounting.models import Wallet
from .in_memory_expense_repository import InMemoryExpenseRepository
from .in_memory_income_repository import InMemoryIncomeRepository
from ..repositories import WalletRepository


class InMemoryWalletRepository(WalletRepository):

    _wallets: dict[str, dict] = {}

    def add_expense(self, wallet, expense):
        wallet = self.get_by_name_with_expenses(wallet.name)
        wallet.add_expense(expense)
        expense.uuid = uuid.uuid4()
        self.add(wallet)
        repository = InMemoryExpenseRepository()
        repository.add(expense)

    def add_income(self, wallet, income):
        wallet = self.get_by_name_with_incomes(wallet.name)
        wallet.add_income(income)
        income.uuid = uuid.uuid4()
        self.add(wallet)
        repository = InMemoryIncomeRepository()
        repository.add(income)

    def add(self, wallet):
        self.__class__._wallets[wallet.name] = {
            'name': wallet.name,
            'balance': wallet.balance,
            'expenses_uuid': [expense.uuid for expense in wallet.expenses],
            'incomes_uuid': [income.uuid for income in wallet.incomes]
        }

    def get_by_name(self, name):
        wallet_data = self._wallets.get(name)
        if wallet_data:
            return Wallet(**wallet_data)

    def get_by_name_with_expenses(self, name):
        wallet_data = self._wallets.get(name)
        expenses = []
        repository = InMemoryExpenseRepository()
        for expense_uuid in wallet_data['expenses_uuid']:
            expense = repository.get_by(uuid=expense_uuid)
            expenses.append(expense)
        return Wallet(**wallet_data, expenses=expenses)

    def get_by_name_with_incomes(self, name):
        wallet_data = self._wallets.get(name)
        incomes = []
        repository = InMemoryIncomeRepository()
        for income_uuid in wallet_data['incomes_uuid']:
            income = repository.get_by_uuid(uuid=income_uuid)
            incomes.append(income)
        return Wallet(**wallet_data, incomes=incomes)
