import datetime
from decimal import Decimal

from pytest import raises

from odin.accounting.controllers.transference_creator import TransferenceCreator
from odin.accounting.repositories import WalletRepository, TransferenceRepository
from tests.factories import WalletBuilder, CategoryFactory


def test_transfer(db_transaction):
    CategoryFactory.create(name='transference')
    wallet_source = WalletBuilder().create()
    wallet_target = WalletBuilder().name('cash').create()
    transference_creator = TransferenceCreator(source=wallet_source, target=wallet_target)
    transference = transference_creator.transfer(amount=Decimal('100_000'))

    wallet_source = WalletRepository().get_by_name(wallet_source.name)
    wallet_target = WalletRepository().get_by_name(wallet_target.name)
    assert wallet_source.balance == Decimal('900_000')
    assert wallet_target.balance == Decimal('1_100_000')
    assert transference.source == wallet_source
    assert transference.target == wallet_target
    assert transference.amount == Decimal('100_000')
    assert transference.date == datetime.date.today()
    assert transference == TransferenceRepository().get_by_uuid(transference.uuid)


def test_transfer_with_date(db_transaction):
    CategoryFactory.create(name='transference')
    wallet_source = WalletBuilder().create()
    wallet_target = WalletBuilder().name('cash').create()
    transference = TransferenceCreator(source=wallet_source, target=wallet_target)
    transference.transfer(amount=Decimal('100_000'), date=datetime.date(2022, 8, 15))
    wallet_source = WalletRepository().get_by_name(wallet_source.name)
    wallet_target = WalletRepository().get_by_name(wallet_target.name)

    assert wallet_source.balance == Decimal('900_000')
    assert wallet_target.balance == Decimal('1_100_000')


def test_transfer_with_date_in_the_future(db_transaction):
    CategoryFactory.create(name='transference')
    wallet_source = WalletBuilder().create()
    wallet_target = WalletBuilder().name('cash').create()
    transference = TransferenceCreator(source=wallet_source, target=wallet_target)
    with raises(ValueError):
        transference.transfer(amount=Decimal('100_000'), date=datetime.date.today() + datetime.timedelta(days=1))
