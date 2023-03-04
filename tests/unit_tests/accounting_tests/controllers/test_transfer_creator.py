import datetime
from decimal import Decimal

from pytest import raises

from odin.accounting.controllers.transfer_creator import TransferCreator
from odin.accounting.repositories import WalletRepository, TransferenceRepository
from tests.factories import WalletBuilder, CategoryFactory


def test_transfer(db_transaction):
    CategoryFactory.create(name='transfer')
    wallet_source = WalletBuilder().create()
    wallet_target = WalletBuilder().name('cash').create()
    transfer_creator = TransferCreator(source=wallet_source, target=wallet_target)
    transfer = transfer_creator.transfer(amount=Decimal('100_000'))

    wallet_source = WalletRepository().get_by_name(wallet_source.name)
    wallet_target = WalletRepository().get_by_name(wallet_target.name)
    assert wallet_source.balance == Decimal('900_000')
    assert wallet_target.balance == Decimal('1_100_000')
    assert transfer.source == wallet_source
    assert transfer.target == wallet_target
    assert transfer.amount == Decimal('100_000')
    assert transfer.date == datetime.date.today()
    assert transfer == TransferenceRepository().get_by_uuid(transfer.uuid)


def test_transfer_with_date(db_transaction):
    CategoryFactory.create(name='transfer')
    wallet_source = WalletBuilder().create()
    wallet_target = WalletBuilder().name('cash').create()
    transfer = TransferCreator(source=wallet_source, target=wallet_target)
    transfer.transfer(amount=Decimal('100_000'), date=datetime.date(2022, 8, 15))
    wallet_source = WalletRepository().get_by_name(wallet_source.name)
    wallet_target = WalletRepository().get_by_name(wallet_target.name)

    assert wallet_source.balance == Decimal('900_000')
    assert wallet_target.balance == Decimal('1_100_000')


def test_transfer_with_date_in_the_future(db_transaction):
    CategoryFactory.create(name='transfer')
    wallet_source = WalletBuilder().create()
    wallet_target = WalletBuilder().name('cash').create()
    transfer = TransferCreator(source=wallet_source, target=wallet_target)
    with raises(ValueError):
        transfer.transfer(amount=Decimal('100_000'), date=datetime.date.today() + datetime.timedelta(days=1))
