import datetime
import re
from decimal import Decimal

from pytest import raises

from odin.controllers import ExpenseCreator, WalletCreator
from tests.factories import CategoryFactory
from tests.utils import UUID_PATTERN


def test_create_expense():
    date = datetime.date.today()
    amount = Decimal('100_000')
    category = CategoryFactory.create()
    wallet_creator = WalletCreator(balance='1_000_000', name='savings account')
    wallet = wallet_creator.create()
    expense_creator = ExpenseCreator(
        date=date,
        amount=amount,
        category=category,
        wallet=wallet
    )
    expense = expense_creator.create()

    assert expense.date == date
    assert expense.amount == amount
    assert re.match(UUID_PATTERN, expense.uuid)
    assert expense.category == category
    assert wallet.balance == Decimal('900_000')


def test_create_expense_with_date_in_the_future():
    wallet_creator = WalletCreator(balance='1_000_000', name='savings account')
    wallet = wallet_creator.create()
    with raises(ValueError) as e:
        ExpenseCreator(
            date=datetime.date.today() + datetime.timedelta(days=2),
            amount='100000',
            category=CategoryFactory.create(),
            wallet=wallet
        )

    assert str(e.value) == 'date must be less or equal than today.'
