import datetime
import re
from decimal import Decimal

from odin.controllers import ExpenseCreator
from odin.repositories import ExpenseRepository
from tests.utils import UUID_PATTERN


def test_create_expense(mocker):
    date = datetime.date(2022, 3, 27)
    amount = Decimal('100_000')
    expense_creator = ExpenseCreator(
        date=date,
        amount=amount
    )
    repository = mocker.patch.object(ExpenseRepository, 'add')
    expense = expense_creator.create()

    assert expense.date == date
    assert expense.amount == amount
    assert re.match(UUID_PATTERN, expense.uuid)
    repository.assert_called_once_with(expense)
