from nyoibo import Entity, fields

from .category import Category


class Expense(Entity):
    _date = fields.DateField()
    _amount = fields.DecimalField()
    _id = fields.StrField(required=True)
    _category = fields.LinkField(to=Category)
