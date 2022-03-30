import datetime
from decimal import Decimal

from odin.controllers import ExpenseCreator, ExpenseGetter


def test_get_expense_by_uuid():
    date = datetime.date(2022, 3, 27)
    amount = Decimal('100_000')
    expense_creator = ExpenseCreator(
        date=date,
        amount=amount
    )
    expense = expense_creator.create()

    expense_getter = ExpenseGetter()
    got_expense = expense_getter.get_by_uuid(uuid=expense.uuid)

    assert got_expense.date == date
    assert got_expense.amount == amount
