import datetime
from decimal import Decimal

from pytest import raises

from odin.repositories import WalletRepository
from tests.factories import WalletBuilder, ExpenseFactory


def test_assert_is_expense_instance(category_fixture):
    wallet = WalletBuilder().build()
    with raises(AssertionError):
        wallet.add_expense(Decimal('100_000'))


def test_add_expense(category_fixture):
    wallet = WalletBuilder().build()
    expense = ExpenseFactory.create(
        date=datetime.date.today(),
        amount='100_000',
        category=category_fixture,
    )
    wallet.add_expense(expense)

    assert wallet.balance == Decimal('900_000')
    assert wallet.expenses == [expense]


def test_get_wallet_with_expenses_from_repository_and_add_expense(category_fixture):
    builder = WalletBuilder().create_expense(amount='100_000')
    wallet = WalletRepository().get_by_name(name=builder.build().name)
    expense = ExpenseFactory.create(
        date=datetime.date.today(),
        amount='100_000',
        category=category_fixture,
    )
    wallet.add_expense(expense)

    assert wallet.balance == Decimal('800_000')
    assert len(wallet.expenses) == 2
