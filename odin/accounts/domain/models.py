from nyoibo import Entity, fields

from .crypto import make_password


class User(Entity):
    _email = fields.StrField()
    _password = fields.StrField(mutable=True)
    _first_name = fields.StrField()
    _last_name = fields.StrField()
    _id = fields.StrField(required=True)

    def __eq__(self, other: 'User'):
        if isinstance(other, User):
            return self._id == other._id
        return False

    def set_password(self, password):
        password = self.__class__._password.parse(password)
        if password:
            # TODO: this is a dummy implementation because is not a priority right now
            if '$' in password:
                self._password = password
            else:
                self._password = make_password(password)

    def check_password(self, raw_password):
        salt = self._password.split('$')[0]
        hash_password = make_password(raw_password, salt=salt)
        return hash_password == self._password


class Token(Entity):
    _value = fields.StrField(required=True)
    _user = fields.LinkField(to=User)
