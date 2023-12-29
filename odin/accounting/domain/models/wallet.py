from nyoibo import Entity, fields

from odin.accounts.domain import User
from .expense import Expense
from .income import Income


class Wallet(Entity):
    _balance = fields.DecimalField(required=True)
    _name = fields.StrField(required=True)
    _expenses: list[Expense] = fields.ListField()
    _incomes: list[Income] = fields.ListField()
    _user = fields.LinkField(to=User)
    _id = fields.StrField(required=True)

    def __init__(self, **kwargs):
        kwargs.setdefault('expenses', [])
        kwargs.setdefault('incomes', [])
        super().__init__(**kwargs)

    def __eq__(self, other: 'Wallet'):
        return self._name == other._name

    def add_expense(self, expense: Expense):
        assert isinstance(expense, Expense), 'expense argument must be Expense instance'
        assert expense.amount <= self._balance, 'expense amount must be lower than wallet balance'
        self._balance -= expense.amount
        self._expenses.append(expense)

    def add_income(self, income: Income):
        assert isinstance(income, Income), 'income argument must be Income instance'
        self._balance += income.amount
        self._incomes.append(income)
