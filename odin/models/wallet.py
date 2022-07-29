import uuid

from nyoibo import Entity, fields

from odin.models import Expense


class Wallet(Entity):
    _balance = fields.DecimalField(required=True)
    _name = fields.StrField(required=True)
    _uuid = fields.StrField(required=True)
    _expenses: list[Expense] = fields.ListField()

    def __init__(self, **kwargs):
        kwargs.setdefault('uuid', uuid.uuid4())
        kwargs.setdefault('expenses', [])
        super().__init__(**kwargs)

    def add_expense(self, expense: Expense):
        assert isinstance(expense, Expense)
        assert expense.amount <= self._balance, 'expense amount must be lower than wallet balance'
        self._balance -= expense.amount
        self._expenses.append(expense)
