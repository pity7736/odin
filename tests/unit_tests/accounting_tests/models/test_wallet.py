import datetime
from decimal import Decimal

from pytest import raises, mark

from odin.accounting.models import Income
from odin.accounting.repositories import WalletRepository
from tests.factories import WalletBuilder, ExpenseFactory


def test_assert_is_expense_instance(category_fixture):
    wallet = WalletBuilder().create()
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
        .add_expense('150_000'),
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
    builder = WalletBuilder().add_expense(amount='100_000')
    repository = WalletRepository()
    wallet = builder.create()
    repository.add(wallet=wallet)
    wallet = repository.get_by_name(name=wallet.name)
    expense = ExpenseFactory.create(
        date=datetime.date.today(),
        amount='100_000',
        category=category_fixture,
    )
    wallet.add_expense(expense)

    assert wallet.balance == Decimal('800_000')
    assert len(wallet.expenses) == 2


def test_add_expense_with_higher_amount_than_wallet_balance(db_transaction):
    wallet = WalletBuilder().balance('100_000').create()
    expense = ExpenseFactory.create(amount=Decimal('100_001'))

    with raises(AssertionError) as error:
        wallet.add_expense(expense)

    assert str(error.value) == 'expense amount must be lower than wallet balance'


def test_add_income(category_fixture):
    wallet = WalletBuilder().create()
    income = Income(
        date=datetime.date.today(),
        amount=Decimal('100_000'),
        category=category_fixture
    )
    wallet.add_income(income)

    assert wallet.balance == Decimal('1_100_000')
    assert len(wallet.incomes) == 1


def test_check_income_type_in_add_income(db_transaction):
    wallet = WalletBuilder().create()
    with raises(AssertionError) as error:
        wallet.add_income(100_000)

    assert str(error.value) == 'income argument must be Income instance'
