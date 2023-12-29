from typing import Optional, Any

from odin.accounting.application.repositories import CategoryRepository, WalletRepository, TransferRepository
from odin.accounting.domain.models import Category, Wallet, Expense, Income, Transfer
from odin.accounts.application.repositories import TokenRepository, UserRepository
from odin.accounts.domain import User


class InMemoryTokenRepository(TokenRepository):

    def __init__(self):
        self._tokens = {}

    def add(self, token):
        self._tokens[token.value] = token

    def get_by_value(self, value):
        return self._tokens.get(value)

    def delete_by_value(self, value):
        self._tokens.pop(value, None)


class InMemoryUserRepository(UserRepository):

    def __init__(self):
        self._users = {}

    def add(self, user):
        self._users[user.email] = User(
            email=user.email,
            password=user.password,
            first_name=user.first_name,
            last_name=user.last_name,
            id=user.id
        )

    def get_by_email(self, email) -> Optional[User]:
        return self._users.get(email)


class InMemoryCategoryRepository(CategoryRepository):

    def __init__(self):
        self._categories: dict[str, dict[str, Any]] = {}

    def add(self, category):
        assert isinstance(category, Category), 'category argument must be Category instance'
        category_name = category.name.lower()
        if category_name in self._categories:
            raise ValueError(f'a category with name {category_name} already exists')
        self._categories[category_name] = {
            'name': category_name,
            'id': category.id
        }

    def get_all(self):
        return tuple(Category(
            name=category_data['name'],
            id=category_data['id']
        ) for category_data in self._categories.values())

    def get_by_name(self, name):
        if name:
            name = name.lower()
            category_data = self._categories.get(name)
            if category_data:
                return Category(name=name, id=category_data['id'])


class InMemoryWalletRepository(WalletRepository):

    def __init__(self):
        self._wallets: dict[str, dict] = {}
        self._expense_repository = InMemoryExpenseRepository()
        self._income_repository = InMemoryIncomeRepository()

    def add_expense(self, wallet, expense):
        wallet = self.get_by_name_with_expenses(wallet.name)
        wallet.add_expense(expense)
        self.add(wallet)
        self._expense_repository.add(expense)

    def add_income(self, wallet, income):
        wallet = self.get_by_name_with_incomes(wallet.name)
        wallet.add_income(income)
        self.add(wallet)
        self._income_repository.add(income)

    def add(self, wallet):
        self._wallets[wallet.name] = {
            'name': wallet.name,
            'balance': wallet.balance,
            'id': wallet.id,
            'expenses_id': [expense.id for expense in wallet.expenses],
            'incomes_id': [income.id for income in wallet.incomes]
        }

    def get_by_name(self, name):
        wallet_data = self._wallets.get(name)
        if wallet_data:
            return Wallet(**wallet_data)

    def get_by_name_with_expenses(self, name):
        wallet_data = self._wallets.get(name)
        expenses = []
        for expense_uuid in wallet_data['expenses_id']:
            expense = self._expense_repository.get_by(uuid=expense_uuid)
            expenses.append(expense)
        return Wallet(**wallet_data, expenses=expenses)

    def get_by_name_with_incomes(self, name):
        wallet_data = self._wallets.get(name)
        incomes = []
        for income_uuid in wallet_data['incomes_id']:
            income = self._income_repository.get_by_id(id=income_uuid)
            incomes.append(income)
        return Wallet(**wallet_data, incomes=incomes)


class InMemoryExpenseRepository:

    def __init__(self):
        self._expenses: dict[str, dict[str, str]] = {}
        self._category_repository = InMemoryCategoryRepository()

    def add(self, expense: Expense):
        self._expenses[expense.id] = {
            'id': expense.id,
            'amount': expense.amount,
            'date': expense.date,
            'category_name': expense.category.name
        }
        self._add_category(expense)

    def _add_category(self, expense):
        try:
            self._category_repository.add(expense.category)
        except ValueError:
            pass

    def get_by(self, uuid) -> Expense:
        try:
            expense_data = self._expenses[uuid]
        except KeyError:
            raise DoesNotExist('Expense not found')
        else:
            return Expense(
                **expense_data,
                category=self._category_repository.get_by_name(expense_data.get('category_name'))
            )


class InMemoryIncomeRepository:

    def __init__(self):
        self._incomes = {}

    def add(self, income: Income):
        self._incomes[income.id] = income

    def get_by_id(self, id: str) -> Income:
        return self._incomes.get(id)


class InMemoryTransferRepository(TransferRepository):

    def __init__(self):
        self._transfers: dict[str, Transfer] = {}

    def add(self, transfer):
        self._transfers[transfer.id] = transfer

    def get_all(self):
        return tuple(self._transfers.values())

    def get_by_id(self, id):
        return self._transfers.get(id)


class DoesNotExist(Exception):
    pass
