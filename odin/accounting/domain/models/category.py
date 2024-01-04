from nyoibo import Entity, fields

from odin.accounts.domain import User


class Category(Entity):
    _id = fields.StrField(required=True)
    _name = fields.StrField(required=True)
    _user = fields.LinkField(to=User)
