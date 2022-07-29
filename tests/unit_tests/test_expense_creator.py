import datetime
import re
from decimal import Decimal

from pytest import raises

from odin.controllers import ExpenseCreator
from tests.factories import CategoryFactory
from tests.utils import UUID_PATTERN


def test_create_expense():
    date = datetime.date.today()
    amount = Decimal('100_000')
    category = CategoryFactory.create()
    expense_creator = ExpenseCreator(
        date=date,
        amount=amount,
        category=category
    )
    expense = expense_creator.create()

    assert expense.date == date
    assert expense.amount == amount
    assert re.match(UUID_PATTERN, expense.uuid)
    assert category == category


def test_create_expense_with_date_in_the_future():
    with raises(ValueError) as e:
        ExpenseCreator(
            date=datetime.date.today() + datetime.timedelta(days=2),
            amount='100000',
            category=CategoryFactory.create()
        )

    assert str(e.value) == 'date must be less or equal than today.'
