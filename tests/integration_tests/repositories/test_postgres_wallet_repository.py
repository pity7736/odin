from decimal import Decimal

from pytest import mark

from tests.factories import WalletBuilder


@mark.asyncio
async def test_get_by_name(db_connection, wallet_repository):
    wallet = await WalletBuilder().create()

    fetched_wallet = await wallet_repository.get_by_name(wallet.name)

    assert fetched_wallet == wallet


@mark.asyncio
async def test_get_by_name_with_expenses(db_connection, wallet_repository):
    wallet = await WalletBuilder().add_expense(Decimal('100_000')).create()

    fetched_wallet = await wallet_repository.get_by_name_with_expenses(wallet.name)

    assert fetched_wallet == wallet
    assert fetched_wallet.expenses == wallet.expenses
    assert fetched_wallet.incomes == []


@mark.asyncio
async def test_get_by_name_with_expenses_when_wallet_has_not_expenses(db_connection, wallet_repository):
    wallet = await WalletBuilder().create()

    fetched_wallet = await wallet_repository.get_by_name_with_expenses(wallet.name)

    assert fetched_wallet == wallet
    assert fetched_wallet.expenses == []


@mark.asyncio
async def test_get_by_name_with_incomes_when_wallet_has_not_incomes(db_connection, wallet_repository):
    wallet = await WalletBuilder().create()

    fetched_wallet = await wallet_repository.get_by_name_with_incomes(wallet.name)

    assert fetched_wallet == wallet
    assert fetched_wallet.incomes == []
    assert fetched_wallet.expenses == []


@mark.asyncio
async def test_get_by_name_with_incomes(db_connection, wallet_repository):
    wallet = await WalletBuilder().add_income(Decimal('100_000')).create()

    fetched_wallet = await wallet_repository.get_by_name_with_expenses(wallet.name)

    assert fetched_wallet == wallet
    assert fetched_wallet.incomes == wallet.expenses
    assert fetched_wallet.expenses == []
