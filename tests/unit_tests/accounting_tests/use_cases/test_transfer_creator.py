import datetime
from decimal import Decimal

from pytest import raises

from odin.accounting.application.use_cases import TransferCreator
from tests.factories import WalletBuilder


def test_transfer(transfer_category, wallet_repository, transfer_repository, category_repository):
    wallet_source = WalletBuilder().create()
    wallet_target = WalletBuilder().name('cash').create()
    transfer_creator = TransferCreator(
        source=wallet_source,
        target=wallet_target,
        wallet_repository=wallet_repository,
        transfer_repository=transfer_repository,
        category_repository=category_repository
    )
    transfer = transfer_creator.transfer(amount=Decimal('100_000'))

    wallet_source = wallet_repository.get_by_name(wallet_source.name)
    wallet_target = wallet_repository.get_by_name(wallet_target.name)
    assert wallet_source.balance == Decimal('900_000')
    assert wallet_target.balance == Decimal('1_100_000')
    assert transfer.source == wallet_source
    assert transfer.target == wallet_target
    assert transfer.amount == Decimal('100_000')
    assert transfer.date == datetime.date.today()
    assert transfer == transfer_repository.get_by_id(transfer.id)


def test_transfer_with_date(transfer_category, wallet_repository, transfer_repository, category_repository):
    wallet_source = WalletBuilder().create()
    wallet_target = WalletBuilder().name('cash').create()
    transfer = TransferCreator(
        source=wallet_source,
        target=wallet_target,
        wallet_repository=wallet_repository,
        transfer_repository=transfer_repository,
        category_repository=category_repository
    )
    transfer.transfer(amount=Decimal('100_000'), date=datetime.date(2022, 8, 15))
    wallet_source = wallet_repository.get_by_name(wallet_source.name)
    wallet_target = wallet_repository.get_by_name(wallet_target.name)

    assert wallet_source.balance == Decimal('900_000')
    assert wallet_target.balance == Decimal('1_100_000')


def test_transfer_with_date_in_the_future(wallet_repository, transfer_category, transfer_repository,
                                          category_repository):
    wallet_source = WalletBuilder().create()
    wallet_target = WalletBuilder().name('cash').create()
    transfer = TransferCreator(
        source=wallet_source,
        target=wallet_target,
        wallet_repository=wallet_repository,
        transfer_repository=transfer_repository,
        category_repository=category_repository
    )
    with raises(ValueError):
        transfer.transfer(amount=Decimal('100_000'), date=datetime.date.today() + datetime.timedelta(days=1))
