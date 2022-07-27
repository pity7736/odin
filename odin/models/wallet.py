import uuid

from nyoibo import Entity, fields


class Wallet(Entity):
    _balance = fields.DecimalField(required=True)
    _name = fields.StrField(required=True)
    _uuid = fields.StrField(required=True)

    def __init__(self, **kwargs):
        kwargs.setdefault('uuid', uuid.uuid4())
        super().__init__(**kwargs)
