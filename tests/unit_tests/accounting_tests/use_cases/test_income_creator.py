import datetime
from decimal import Decimal

from odin.accounting.application.use_cases import IncomeCreator
from tests.factories import WalletBuilder


def test_create_income(category_fixture, wallet_repository):
    date = datetime.date.today()
    amount = Decimal('100_000')
    wallet = WalletBuilder().create()
    income_creator = IncomeCreator(
        date=date,
        amount=amount,
        category=category_fixture,
        wallet=wallet,
        wallet_repository=wallet_repository
    )
    income = income_creator.create()
    wallet = wallet_repository.get_by_name(wallet.name)

    assert income.date == date
    assert income.amount == amount
    assert wallet.balance == Decimal('1_100_000')
