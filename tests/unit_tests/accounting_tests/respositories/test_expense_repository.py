import uuid

from pytest import raises

from odin.accounting.models import Expense, Category
from odin.accounting.repositories.in_memory_reposotiries import InMemoryExpenseRepository
from odin.accounting.repositories.exceptions import DoesNotExist
from tests.factories import ExpenseFactory


def test_get_expense_by_uuid(expense_fixture):
    repository = InMemoryExpenseRepository()
    gotten_expense = repository.get_by(uuid=expense_fixture.uuid)

    assert gotten_expense.date == expense_fixture.date
    assert gotten_expense.amount == expense_fixture.amount
    assert gotten_expense.uuid == expense_fixture.uuid
    assert gotten_expense is not expense_fixture


def test_get_expense_by_uuid_twice(expense_fixture):
    repository = InMemoryExpenseRepository()
    repository.get_by(uuid=expense_fixture.uuid)
    gotten_expense = repository.get_by(uuid=expense_fixture.uuid)

    assert gotten_expense.date == expense_fixture.date
    assert gotten_expense.amount == expense_fixture.amount
    assert gotten_expense.uuid == expense_fixture.uuid
    assert gotten_expense is not expense_fixture


def test_get_non_existing_expense_by_uuid(expense_fixture):
    repository = InMemoryExpenseRepository()
    with raises(DoesNotExist) as e:
        repository.get_by(uuid=uuid.uuid4())

    assert str(e.value) == 'Expense not found'


def test_get_all_expenses(db_transaction):
    ExpenseFactory.create_batch(5)
    repository = InMemoryExpenseRepository()
    expenses = repository.get_all()

    assert len(expenses) == 5
    for expense in expenses:
        assert isinstance(expense, Expense)
        assert isinstance(expense.category, Category)


def test_get_all_expenses_twice(db_transaction):
    ExpenseFactory.create_batch(5)
    repository = InMemoryExpenseRepository()
    repository.get_all()
    expenses = repository.get_all()

    assert len(expenses) == 5
    for expense in expenses:
        assert isinstance(expense, Expense)
        assert isinstance(expense.category, Category)
