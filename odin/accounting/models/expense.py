from nyoibo import Entity, fields

from .category import Category


class Expense(Entity):
    _date = fields.DateField()
    _amount = fields.DecimalField()
    _uuid = fields.StrField(mutable=True)
    _category = fields.LinkField(to=Category)

    def __init__(self, **kwargs):
        super().__init__(**kwargs)
