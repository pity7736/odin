import datetime
from decimal import Decimal

from pytest import mark

from odin.accounting.application.use_cases import IncomeCreator
from tests.factories import WalletBuilder


@mark.asyncio
async def test_create_income(category_fixture, wallet_repository):
    date = datetime.date.today()
    amount = Decimal('100_000')
    wallet = await WalletBuilder().create()
    income_creator = IncomeCreator(
        date=date,
        amount=amount,
        category=category_fixture,
        wallet=wallet,
        wallet_repository=wallet_repository
    )
    income = await income_creator.create()

    assert income.date == date
    assert income.amount == amount
    assert wallet.balance == Decimal('1_100_000')
