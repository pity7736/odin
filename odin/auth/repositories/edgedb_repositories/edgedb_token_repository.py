from odin.accounting.repositories.edgedb_repositories.db_client import DBClient
from odin.accounts.models import User
from odin.auth.models import Token
from ..repositories import TokenRepository


class EdgeDBTokenRepository(TokenRepository):

    def __init__(self):
        self._client = DBClient()

    def add(self, token):
        self._client.execute(
            'insert Token {value := <str>$value, user := (select User {id} filter User.email = <str>$email)}',
            value=token.value,
            email=token.user.email
        )

    def get_by_value(self, value):
        record = self._client.query_single(
            'select Token {value, user: {email, password}} filter .value = <str> $value',
            value=value
        )
        if record:
            return Token(
                value=record.value,
                user=User(
                    email=record.user.email,
                    password=record.user.password
                )
            )

    def delete_by_value(self, value):
        self._client.execute(
            'delete Token filter .value = <str>$value',
            value=value
        )
