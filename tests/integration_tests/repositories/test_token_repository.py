from odin.auth.models import Token
from odin.auth.repositories.edgedb_repositories import EdgeDBTokenRepository
from odin.utils import get_random_string


def test_get_by_value(user_fixture):
    token = Token(
        user=user_fixture,
        value=get_random_string()
    )
    repository = EdgeDBTokenRepository()
    repository.add(token)

    fetched_token = repository.get_by_value(value=token.value)

    assert fetched_token.value == token.value
    assert fetched_token.user.email == user_fixture.email
