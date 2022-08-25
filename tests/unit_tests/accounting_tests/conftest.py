from pytest import fixture
from starlette.testclient import TestClient

from odin.accounting.api import app
from tests.factories import ExpenseFactory, CategoryFactory


@fixture
def test_client():
    return TestClient(app=app)


@fixture
def expense_fixture(db_transaction):
    return ExpenseFactory.create(
        date='2022-03-27',
        amount='100_00'
    )


@fixture
def category_fixture(db_transaction):
    return CategoryFactory.create()
