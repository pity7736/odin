import datetime
from decimal import Decimal

from pytest import raises, mark

from odin.accounting.application.use_cases import TransferCreator
from tests.factories import WalletBuilder


@mark.asyncio
async def test_transfer(transfer_category, wallet_repository, transfer_repository, category_repository,
                        user_fixture):
    category_repository.get_by_name.return_value = transfer_category
    wallet_source = await WalletBuilder().user(user_fixture).create()
    wallet_target = await WalletBuilder().user(user_fixture).name('cash').create()
    transfer_creator = TransferCreator(
        source=wallet_source,
        target=wallet_target,
        wallet_repository=wallet_repository,
        transfer_repository=transfer_repository,
        category_repository=category_repository
    )
    transfer = await transfer_creator.transfer(amount=Decimal('100_000'))

    assert wallet_source.balance == Decimal('900_000')
    assert wallet_target.balance == Decimal('1_100_000')
    assert transfer.source == wallet_source
    assert transfer.target == wallet_target
    assert transfer.amount == Decimal('100_000')
    assert transfer.date == datetime.date.today()


@mark.asyncio
async def test_transfer_with_date(transfer_category, wallet_repository, transfer_repository, category_repository,
                                  user_fixture):
    category_repository.get_by_name.return_value = transfer_category
    wallet_source = await WalletBuilder().user(user_fixture).create()
    wallet_target = await WalletBuilder().user(user_fixture).name('cash').create()
    transfer = TransferCreator(
        source=wallet_source,
        target=wallet_target,
        wallet_repository=wallet_repository,
        transfer_repository=transfer_repository,
        category_repository=category_repository
    )
    await transfer.transfer(amount=Decimal('100_000'), date=datetime.date(2022, 8, 15))

    assert wallet_source.balance == Decimal('900_000')
    assert wallet_target.balance == Decimal('1_100_000')


@mark.asyncio
async def test_transfer_with_date_in_the_future(wallet_repository, transfer_category, transfer_repository,
                                                category_repository):
    category_repository.get_by_name.return_value = transfer_category
    wallet_source = await WalletBuilder().create()
    wallet_target = await WalletBuilder().name('cash').create()
    transfer = TransferCreator(
        source=wallet_source,
        target=wallet_target,
        wallet_repository=wallet_repository,
        transfer_repository=transfer_repository,
        category_repository=category_repository
    )
    with raises(ValueError):
        await transfer.transfer(amount=Decimal('100_000'), date=datetime.date.today() + datetime.timedelta(days=1))
