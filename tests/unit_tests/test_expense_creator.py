import datetime
import re
from decimal import Decimal

from pytest import raises

from odin.controllers import ExpenseCreator
from tests.utils import UUID_PATTERN


def test_create_expense():
    date = datetime.date.today()
    amount = Decimal('100_000')
    expense_creator = ExpenseCreator(
        date=date,
        amount=amount
    )
    expense = expense_creator.create()

    assert expense.date == date
    assert expense.amount == amount
    assert re.match(UUID_PATTERN, expense.uuid)


def test_create_expense_with_date_in_the_future():
    with raises(ValueError) as e:
        ExpenseCreator(
            date=datetime.date.today() + datetime.timedelta(days=2),
            amount='100000'
        )

    assert str(e.value) == 'date must be less or equal than today.'
