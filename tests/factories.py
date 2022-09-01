import datetime
from decimal import Decimal
from uuid import uuid4

import factory

from odin.accounting.controllers import WalletCreator, ExpenseCreator
from odin.accounting.models import Expense, Category, Wallet
from odin.accounting.repositories import ExpenseRepository, CategoryRepository


class CategoryFactory(factory.Factory):
    name = factory.Sequence(lambda n: f'test category{n}')

    class Meta:
        model = Category

    @classmethod
    def _create(cls, model_class, *args, **kwargs):
        category = super()._create(model_class, *args, **kwargs)
        repository = CategoryRepository()
        repository.add(category)
        return category


class ExpenseFactory(factory.Factory):
    uuid = factory.LazyFunction(uuid4)
    date = datetime.date(2022, 3, 30)
    amount = Decimal('100_000')
    category = factory.SubFactory(CategoryFactory)

    class Meta:
        model = Expense

    @classmethod
    def _create(cls, model_class, *args, **kwargs):
        expense = super()._create(model_class, *args, **kwargs)
        repository = ExpenseRepository()
        repository.add(expense)
        return expense


class WalletBuilder:

    def __init__(self):
        self._name = 'savings account'
        self._balance = '1_000_000'
        self._expenses_data = []

    def name(self, name) -> 'WalletBuilder':
        self._name = name
        return self

    def balance(self, balance) -> 'WalletBuilder':
        self._balance = balance
        return self

    def add_expense(self, amount, date=None, category=None) -> 'WalletBuilder':
        self._expenses_data.append({
            'amount': amount,
            'date': date or datetime.date.today(),
            'category': category
        })
        return self

    def create(self) -> Wallet:
        wallet = WalletCreator(name=self._name, balance=self._balance).create()
        for expense_data in self._expenses_data:
            ExpenseCreator(
                amount=expense_data['amount'],
                date=expense_data['date'],
                category=expense_data['category'] or CategoryFactory.create(),
                wallet=wallet
            ).create()
        return wallet

    def build(self) -> Wallet:
        wallet = Wallet(name=self._name, balance=self._balance)
        for expense_data in self._expenses_data:
            expense = Expense(
                amount=expense_data['amount'],
                date=expense_data['date'],
                category=expense_data['category'] or CategoryFactory.build(),
                wallet=wallet
            )
            wallet.add_expense(expense)
        return wallet
