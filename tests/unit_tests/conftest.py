import uuid
from unittest.mock import AsyncMock

from pytest import fixture
from pytest_asyncio import fixture as async_fixture

from odin.accounting.infrastructure.repositories import RepositoryFactory
from odin.accounting.infrastructure.repositories.postgres_repositories import PostgresCategoryRepository, \
    PostgresWalletRepository, PostgresTransferRepository
from odin.accounts.domain.models import User, Token
from odin.accounts.infrastructure.repositories.postgres_repositories import PostgresUserRepository, \
    PostgresTokenRepository
from odin.accounts.domain.crypto import get_random_string
from tests.repositories import InMemoryTokenRepository, InMemoryUserRepository
from tests.factories import CategoryFactory


@fixture(autouse=True)
def user_repository(mocker):
    repository = InMemoryUserRepository()
    mocker.patch.object(PostgresUserRepository, '__new__', return_value=repository)
    return repository


@fixture(autouse=True)
def token_repository(mocker):
    repository = InMemoryTokenRepository()
    mocker.patch.object(PostgresTokenRepository, '__new__', return_value=repository)
    return repository


@fixture(autouse=True)
def category_repository(mocker):
    repository = AsyncMock(spec_set=PostgresCategoryRepository)
    mocker.patch.object(RepositoryFactory, 'get_category_repository', return_value=repository)
    return repository


@fixture(autouse=True)
def wallet_repository(mocker):
    repository = AsyncMock(spec_set=PostgresWalletRepository)
    mocker.patch.object(RepositoryFactory, 'get_wallet_repository', return_value=repository)
    return repository


@fixture(autouse=True)
def transfer_repository(mocker):
    repository = AsyncMock(spec_set=PostgresTransferRepository)
    mocker.patch.object(RepositoryFactory, 'get_transfer_repository', return_value=repository)
    return repository


@async_fixture
async def user_fixture(user_repository):
    user = User(
        email='me@raiseexception.com',
        password='test',
        first_name='julián',
        last_name='cortés',
        id=uuid.uuid4()
    )
    await user_repository.add(user)
    return user


@async_fixture
async def token_value_fixture(user_fixture, test_client, token_repository):
    token = Token(
        value=get_random_string(length=50),
        user=user_fixture
    )
    await token_repository.add(token)
    return token.value


@async_fixture
async def category_fixture(category_repository, user_fixture):
    return await CategoryFactory.create(user=user_fixture)


@async_fixture
async def transfer_category(category_repository, user_fixture):
    return await CategoryFactory.create(name='transfer', user=user_fixture)
