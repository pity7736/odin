from decimal import Decimal

from odin.models import Wallet
from odin.repositories import WalletRepository


def test_get_by_name(db_transaction):
    wallet_name = 'savings account'
    repository = WalletRepository()
    repository.add(wallet=Wallet(name=wallet_name, balance='100_000'))
    wallet = repository.get_by_name(wallet_name)

    assert wallet.name == wallet_name
    assert wallet.balance == Decimal('100_000')
