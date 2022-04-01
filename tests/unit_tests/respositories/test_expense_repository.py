import uuid

from pytest import raises

from odin.repositories import ExpenseRepository
from odin.repositories.exceptions import DoesNotExist


def test_get_expense_by_uuid(expense_fixture):
    repository = ExpenseRepository()
    gotten_expense = repository.get_by(uuid=expense_fixture.uuid)

    assert gotten_expense.date == expense_fixture.date
    assert gotten_expense.amount == expense_fixture.amount
    assert gotten_expense.uuid == expense_fixture.uuid


def test_get_non_existing_expense_by_uuid(expense_fixture):
    repository = ExpenseRepository()
    with raises(DoesNotExist) as e:
        repository.get_by(uuid=uuid.uuid4())

    assert str(e.value) == 'Expense not found'
