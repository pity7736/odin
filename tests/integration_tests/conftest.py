from pytest import fixture

from odin.accounts.models import User
from odin.accounts.repositories.edgedb_repositories import EdgeDBUserRepository


@fixture
def user_fixture(db_transaction):
    user = User(
        email='me@raiseexception.com',
        password='test',
        first_name='julián',
        last_name='cortés'
    )
    EdgeDBUserRepository().add(user)
    return user
