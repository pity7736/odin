from nyoibo import Entity, fields

from odin.accounts.models import User


class Token(Entity):
    _value = fields.StrField()
    _user = fields.LinkField(to=User)
