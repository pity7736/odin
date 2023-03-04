import uuid

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
    _uuid = fields.StrField()

    def __init__(self, **kwargs):
        kwargs.setdefault('uuid', uuid.uuid4())
        super().__init__(**kwargs)
