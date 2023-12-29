from nyoibo import Entity, fields


class Category(Entity):
    _id = fields.StrField(required=True)
    _name = fields.StrField(required=True)
