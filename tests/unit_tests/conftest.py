from pytest import fixture

from odin.accounting.repositories import ExpenseRepository, WalletRepository, CategoryRepository, TransferenceRepository
from odin.accounts.models import User
from odin.accounts.repositories import InMemoryUserRepository
from odin.auth.models import Token
from odin.auth.repositories.in_memory_repositories import InMemoryTokenRepository
from odin.utils import get_random_string


@fixture
def db_transaction():
    yield
    ExpenseRepository._expenses.clear()
    WalletRepository._wallets.clear()
    CategoryRepository._categories.clear()
    TransferenceRepository._transfers.clear()
    InMemoryUserRepository._user.clear()
    InMemoryTokenRepository._tokens.clear()


@fixture
def user_fixture(db_transaction):
    user = User(
        email='me@raiseexception.com',
        password='test',
        first_name='julián',
        last_name='cortés'
    )
    InMemoryUserRepository().add(user)
    return user


@fixture
def token_value_fixture(user_fixture, test_client):
    token = Token(
        value=get_random_string(length=50),
        user=user_fixture
    )
    repository = InMemoryTokenRepository()
    repository.add(token)
    return token.value
