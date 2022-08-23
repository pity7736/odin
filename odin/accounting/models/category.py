from nyoibo import Entity, fields


class Category(Entity):
    _name = fields.StrField(required=True)
