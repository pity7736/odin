import uuid

from nyoibo import Entity, fields

from .category import Category


class Expense(Entity):
    _date = fields.DateField()
    _amount = fields.DecimalField()
    _uuid = fields.StrField()
    _category = fields.LinkField(to=Category)

    def __init__(self, **kwargs):
        kwargs.setdefault('uuid', uuid.uuid4())
        super().__init__(**kwargs)
