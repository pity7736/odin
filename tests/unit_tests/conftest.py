from pytest import fixture

from odin.accounting.repositories.in_memory_reposotiries import InMemoryWalletRepository, InMemoryCategoryRepository, \
    InMemoryTransferRepository
from odin.accounting.repositories.in_memory_reposotiries.in_memory_expense_repository import InMemoryExpenseRepository
from odin.accounting.repositories.in_memory_reposotiries.in_memory_income_repository import InMemoryIncomeRepository
from odin.accounts.domain.models import User, Token
from odin.accounts.infrastructure.repositories.postgres_repositories import PostgresUserRepository, \
    PostgresTokenRepository
from odin.accounts.domain.crypto import get_random_string
from tests.repositories import InMemoryTokenRepository, InMemoryUserRepository


@fixture
def db_transaction():
    yield
    InMemoryExpenseRepository._expenses.clear()
    InMemoryIncomeRepository._incomes.clear()
    InMemoryWalletRepository._wallets.clear()
    InMemoryCategoryRepository._categories.clear()
    InMemoryTransferRepository._transfers.clear()


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
def user_fixture(db_transaction, user_repository):
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
