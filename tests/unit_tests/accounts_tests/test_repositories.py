from odin.accounts.repositories import InMemoryUserRepository

from odin.accounts.models import User


def test_get_by_email(db_transaction):
    user = User(
        email='me@raiseexception.com',
        password='test',
        first_name='julián',
        last_name='cortés'
    )
    repository = InMemoryUserRepository()
    repository.add(user)

    fetched_user = repository.get_by_email(user.email)

    assert user.password == fetched_user.password
    assert user.email == fetched_user.email
