from nyoibo import Entity, fields

from .category import Category


class Expense(Entity):
    _date = fields.DateField()
    _amount = fields.DecimalField()
    _id = fields.StrField(required=True)
    _category = fields.LinkField(to=Category)

    def __eq__(self, other: 'Expense'):
        if isinstance(other, Expense):
            return self._id == other.id
        return False
