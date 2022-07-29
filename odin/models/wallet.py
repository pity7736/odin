import uuid

from nyoibo import Entity, fields

from odin.models import Expense


class Wallet(Entity):
    _balance = fields.DecimalField(required=True)
    _name = fields.StrField(required=True)
    _uuid = fields.StrField(required=True)
    _expenses = fields.ListField()

    def __init__(self, **kwargs):
        kwargs.setdefault('uuid', uuid.uuid4())
        super().__init__(**kwargs)

    def add_expense(self, expense: Expense):
        self._balance -= expense.amount
