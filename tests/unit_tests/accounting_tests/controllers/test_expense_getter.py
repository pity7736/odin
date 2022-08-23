import uuid

from pytest import fixture

from odin.accounting.controllers import ExpenseGetter
from tests.factories import ExpenseFactory


@fixture(autouse=True)
def auto_db_transaction(db_transaction):
    pass


def test_get_expense_by_uuid(expense_fixture):
    expense_getter = ExpenseGetter()
    gotten_expense = expense_getter.get_by_uuid(uuid=expense_fixture.uuid)

    assert gotten_expense.date == expense_fixture.date
    assert gotten_expense.amount == expense_fixture.amount


def test_get_non_existing_expense_by_uuid(expense_fixture):
    expense_fixture = ExpenseGetter()
    gotten_expense = expense_fixture.get_by_uuid(uuid=uuid.uuid4())

    assert gotten_expense is None


def test_get_all(db_transaction):
    ExpenseFactory.create_batch(2)
    expense_getter = ExpenseGetter()
    expenses = expense_getter.all()

    assert expenses == expenses
