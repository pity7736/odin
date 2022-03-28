from nyoibo import Entity, fields


class Expense(Entity):
    _date = fields.DateField()
    _amount = fields.DecimalField()
    _uuid = fields.StrField()
