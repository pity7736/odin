import uuid

from odin.controllers import ExpenseGetter, ExpenseCreator
from odin.repositories import ExpenseRepository


def test_get_expense_by_uuid(expense_fixture):
    expense_getter = ExpenseGetter()
    gotten_expense = expense_getter.get_by_uuid(uuid=expense_fixture.uuid)

    assert gotten_expense.date == expense_fixture.date
    assert gotten_expense.amount == expense_fixture.amount


def test_get_non_existing_expense_by_uuid(expense_fixture):
    expense_fixture = ExpenseGetter()
    gotten_expense = expense_fixture.get_by_uuid(uuid=uuid.uuid4())

    assert gotten_expense is None


def test_get_all(mocker):
    expense_creator = ExpenseCreator(
        date='2022-03-30',
        amount='100'
    )
    expense0 = expense_creator.create()
    expense_creator = ExpenseCreator(
        date='2022-03-29',
        amount='10000'
    )
    expense1 = expense_creator.create()

    expense_getter = ExpenseGetter()
    mocker.patch.object(ExpenseRepository, 'get_all', return_value=(expense0, expense1))
    expenses = expense_getter.all()

    assert expenses == (expense0, expense1)
