import datetime
from decimal import Decimal

from odin.controllers import ExpenseCreator, ExpenseGetter
from odin.repositories import ExpenseRepository


def test_get_expense_by_uuid(mocker):
    date = datetime.date(2022, 3, 27)
    amount = Decimal('100_000')
    expense_creator = ExpenseCreator(
        date=date,
        amount=amount
    )
    expense = expense_creator.create()

    repository = mocker.patch.object(ExpenseRepository, 'get_by', return_value=expense)
    expense_getter = ExpenseGetter()
    got_expense = expense_getter.get_by_uuid(uuid=expense.uuid)

    repository.assert_called_once_with(uuid=expense.uuid)
    assert got_expense.date == date
    assert got_expense.amount == amount
