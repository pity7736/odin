from pytest import fixture

from odin.accounting.infrastructure.repositories.postgres_repositories import PostgresCategoryRepository, \
    PostgresWalletRepository, PostgresTransferRepository
from odin.accounts.domain.models import User, Token
from odin.accounts.infrastructure.repositories.postgres_repositories import PostgresUserRepository, \
    PostgresTokenRepository
from odin.accounts.domain.crypto import get_random_string
from tests.repositories import InMemoryTokenRepository, InMemoryUserRepository, InMemoryCategoryRepository, \
    InMemoryWalletRepository, InMemoryTransferRepository
from tests.factories import CategoryFactory


@fixture
def user_repository(mocker):
    repo = InMemoryUserRepository()
    mocker.patch.object(PostgresUserRepository, '__new__', return_value=repo)
    return repo


@fixture
def token_repository(mocker):
    repo = InMemoryTokenRepository()
    mocker.patch.object(PostgresTokenRepository, '__new__', return_value=repo)
    return repo


@fixture
def category_repository(mocker):
    repo = InMemoryCategoryRepository()
    mocker.patch.object(PostgresCategoryRepository, '__new__', return_value=repo)
    return repo


@fixture
def wallet_repository(mocker):
    repo = InMemoryWalletRepository()
    mocker.patch.object(PostgresWalletRepository, '__new__', return_value=repo)
    return repo


@fixture
def transfer_repository(mocker):
    repo = InMemoryTransferRepository()
    mocker.patch.object(PostgresTransferRepository, '__new__', return_value=repo)
    return repo


@fixture
def user_fixture(user_repository):
    user = User(
        email='me@raiseexception.com',
        password='test',
        first_name='julián',
        last_name='cortés'
    )
    user_repository.add(user)
    return user


@fixture
def token_value_fixture(user_fixture, test_client, token_repository):
    token = Token(
        value=get_random_string(length=50),
        user=user_fixture
    )
    token_repository.add(token)
    return token.value


@fixture
def category_fixture(category_repository):
    return CategoryFactory.create()


@fixture
def transfer_category(category_repository):
    return CategoryFactory.create(name='transfer')
