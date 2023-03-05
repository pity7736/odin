import datetime
from decimal import Decimal

from odin.accounting.controllers import IncomeCreator
from odin.accounting.repositories.repository_factory import get_wallet_repository
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
    wallet = get_wallet_repository().get_by_name(wallet.name)

    assert income.date == date
    assert income.amount == amount
    assert wallet.balance == Decimal('1_100_000')
