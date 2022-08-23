import datetime
from decimal import Decimal

from odin.controllers import IncomeCreator
from odin.repositories import WalletRepository
from tests.factories import WalletBuilder


def test_create_income(category_fixture):
    date = datetime.date.today()
    amount = Decimal('100_000')
    wallet = WalletBuilder().create()
    income_creator = IncomeCreator(
        date=date,
        amount=amount,
        category=category_fixture,
        wallet=wallet
    )
    income = income_creator.create()
    wallet = WalletRepository().get_by_name(wallet.name)

    assert income.date == date
    assert income.amount == amount
    assert wallet.balance == Decimal('1_100_000')
