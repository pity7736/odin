from odin.accounts.models import User
from odin.accounts.repositories.edgedb_repositories import EdgeDBUserRepository


def test_get_by_email(db_transaction):
    repository = EdgeDBUserRepository()
    user = User(
        email='email@user.com',
        password='password',
        first_name='john',
        last_name='doe'
    )
    repository.add(user)

    fetched_user = repository.get_by_email(email=user.email)

    assert fetched_user.email == user.email
