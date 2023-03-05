from decimal import Decimal

from odin.accounting.models import Wallet
from odin.accounting.repositories.repository_factory import get_wallet_repository
from tests.factories import WalletBuilder


def test_get_by_name(db_transaction):
    wallet_name = 'savings account'
    repository = get_wallet_repository()
    repository.add(wallet=Wallet(name=wallet_name, balance='100_000'))
    wallet = repository.get_by_name(wallet_name)

    assert wallet.name == wallet_name
    assert wallet.balance == Decimal('100_000')
    assert len(wallet.expenses) == 0


def test_get_by_wrong_name(db_transaction):
    WalletBuilder().create()
    repository = get_wallet_repository()
    assert repository.get_by_name('wrong name') is None
