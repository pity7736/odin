from nyoibo import Entity, fields

from odin.accounts.domain import User
from ..enumerations import CategoryType


class Category(Entity):
    _id = fields.StrField(required=True)
    _name = fields.StrField(required=True)
    _user = fields.LinkField(to=User)
    _type = fields.StrField(required=True, choices=CategoryType)
