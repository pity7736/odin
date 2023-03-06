from odin.utils import DBClient
from odin.accounts.models import User
from ..repositories import UserRepository


class EdgeDBUserRepository(UserRepository):

    def __init__(self):
        self._client = DBClient()

    def add(self, user):
        self._client.execute(
            'insert User {'
            'email := <str>$email,'
            'password := <str>$password,'
            'first_name := <str>$first_name,'
            'last_name := <str>$last_name}',
            email=user.email,
            password=user.password,
            first_name=user.first_name,
            last_name=user.last_name
        )

    def get_by_email(self, email):
        record = self._client.query_single('select User {email, password} filter .email = <str>$email', email=email)
        return User(
            email=record.email,
            password=record.password
        )
