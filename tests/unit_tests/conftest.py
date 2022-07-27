from pytest import fixture
from starlette.testclient import TestClient

from odin.api import app
from odin.repositories import ExpenseRepository, CategoryRepository, WalletRepository
from tests.factories import ExpenseFactory


@fixture
def expense_fixture():
    return ExpenseFactory.create(
        date='2022-03-27',
        amount='100_00'
    )


@fixture
def test_client():
    return TestClient(app=app)


@fixture
def db_transaction():
    yield
    CategoryRepository._categories.clear()
    ExpenseRepository._expenses.clear()
    WalletRepository._wallets.clear()
