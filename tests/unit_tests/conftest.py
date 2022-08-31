from _pytest.fixtures import fixture
from starlette.testclient import TestClient

from odin.accounts.models import User
from odin.auth.repositories import TokenRepository
from odin.main import app
from odin.accounting.repositories import ExpenseRepository, WalletRepository, CategoryRepository, TransferenceRepository
from odin.accounts.repositories import UserRepository


@fixture
def db_transaction():
    yield
    ExpenseRepository._expenses.clear()
    WalletRepository._wallets.clear()
    CategoryRepository._categories.clear()
    TransferenceRepository._transfers.clear()
    UserRepository._user.clear()
    TokenRepository._tokens.clear()


@fixture
def test_client():
    return TestClient(app=app)


@fixture
def user_fixture(db_transaction):
    user = User(
        email='me@raiseexception.com',
        password='test',
        first_name='julián',
        last_name='cortés'
    )
    UserRepository().add(user)
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
