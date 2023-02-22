from pytest import fixture

from odin.accounting.repositories import ExpenseRepository, WalletRepository, CategoryRepository, TransferenceRepository
from odin.accounts.models import User
from odin.accounts.repositories import InMemoryUserRepository
from odin.auth.repositories.in_memory_repositories import InMemoryTokenRepository


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
    login_response = test_client.post(
        '/auth/login',
        json={
            'email': user_fixture.email,
            'password': 'test'
        }
    )
    data = login_response.json()
    return data['token']
