import datetime
from decimal import Decimal

from pytest import raises, mark

from odin.repositories import WalletRepository
from tests.factories import WalletBuilder, ExpenseFactory


def test_assert_is_expense_instance(category_fixture):
    wallet = WalletBuilder().build()
    with raises(AssertionError) as error:
        wallet.add_expense(Decimal('100_000'))

    assert str(error.value) == 'expense argument must be Expense instance'


add_expense_params = (
    (
        WalletBuilder(),
        Decimal('100_000'),
        Decimal('900_000'),
        1
    ),
    (
        WalletBuilder()
        .create_expense('150_000'),
        Decimal('200_000'),
        Decimal('650_000'),
        2
    )
)


@mark.parametrize('wallet_builder, amount, expected_balance, expected_expenses_number', add_expense_params)
def test_add_expense(wallet_builder, amount, expected_balance, expected_expenses_number, category_fixture):
    wallet = wallet_builder.build()
    expense = ExpenseFactory.create(
        date=datetime.date.today(),
        amount=amount,
        category=category_fixture,
    )
    wallet.add_expense(expense)

    assert wallet.balance == expected_balance
    assert len(wallet.expenses) == expected_expenses_number


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


def test_add_expense_with_higher_amount_than_wallet_balance(db_transaction):
    wallet = WalletBuilder().balance('100_000').build()
    expense = ExpenseFactory.create(amount=Decimal('100_001'))

    with raises(AssertionError) as error:
        wallet.add_expense(expense)

    assert str(error.value) == 'expense amount must be lower than wallet balance'
