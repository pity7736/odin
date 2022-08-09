import datetime
import re
from decimal import Decimal

from pytest import raises

from odin.controllers import ExpenseCreator
from odin.repositories import WalletRepository
from tests.factories import WalletBuilder
from tests.utils import UUID_PATTERN


def test_create_expense(category_fixture):
    date = datetime.date.today()
    amount = Decimal('100_000')
    wallet = WalletBuilder().build()
    expense_creator = ExpenseCreator(
        date=date,
        amount=amount,
        category=category_fixture,
        wallet=wallet
    )
    expense = expense_creator.create()
    wallet = WalletRepository().get_by_name(wallet.name)

    assert expense.date == date
    assert expense.amount == amount
    assert re.match(UUID_PATTERN, expense.uuid)
    assert expense.category == category_fixture
    assert wallet.balance == Decimal('900_000')


def test_create_expense_with_date_in_the_future(category_fixture):
    wallet = WalletBuilder().build()
    with raises(ValueError) as e:
        ExpenseCreator(
            date=datetime.date.today() + datetime.timedelta(days=2),
            amount='100000',
            category=category_fixture,
            wallet=wallet
        )

    assert str(e.value) == 'date must be less or equal than today.'


def test_without_category(db_transaction):
    wallet = WalletBuilder().build()
    with raises(ValueError):
        ExpenseCreator(
            date=datetime.date.today(),
            amount=Decimal('100_000'),
            wallet=wallet
        )


def test_without_wallet(category_fixture):
    with raises(ValueError):
        ExpenseCreator(
            date=datetime.date.today(),
            amount=Decimal('100_000'),
            category=category_fixture
        )
