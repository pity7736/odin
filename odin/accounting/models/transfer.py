from nyoibo import Entity, fields

from .expense import Expense
from .income import Income
from .wallet import Wallet


class Transfer(Entity):
    _source = fields.LinkField(to=Wallet)
    _target = fields.LinkField(to=Wallet)
    _expense = fields.LinkField(to=Expense)
    _income = fields.LinkField(to=Income)
    _amount = fields.DecimalField()
    _date = fields.DateField()
    _uuid = fields.StrField(mutable=True)

    def __init__(self, **kwargs):
        super().__init__(**kwargs)
