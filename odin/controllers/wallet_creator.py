from nyoibo import Entity, fields

from odin.models import Wallet
from odin.repositories import WalletRepository


class WalletCreator(Entity):
    _balance = fields.DecimalField(required=True, private=True)
    _name = fields.StrField(required=True, private=True)

    def create(self) -> Wallet:
        wallet = Wallet(
            name=self._name,
            balance=self._balance
        )
        repository = WalletRepository()
        repository.add(wallet)
        return wallet
