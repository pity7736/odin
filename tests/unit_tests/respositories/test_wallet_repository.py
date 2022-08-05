from decimal import Decimal

from odin.models import Wallet
from odin.repositories import WalletRepository
from tests.factories import WalletBuilder


def test_get_by_name(db_transaction):
    wallet_name = 'savings account'
    repository = WalletRepository()
    repository.add(wallet=Wallet(name=wallet_name, balance='100_000'))
    wallet = repository.get_by_name(wallet_name)

    assert wallet.name == wallet_name
    assert wallet.balance == Decimal('100_000')
    assert len(wallet.expenses) == 0


def test_get_by_wrong_name(db_transaction):
    WalletBuilder().build()
    repository = WalletRepository()
    assert repository.get_by_name('wrong name') is None


def test_update_non_existing_wallet(category_fixture):
    wallet_name = 'savings account'
    repository = WalletRepository()
    repository.update(Wallet(name=wallet_name, balance='100000'))

    assert repository.get_by_name(wallet_name) is None
