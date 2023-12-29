from nyoibo import Entity, fields

from ..repositories import WalletRepository
from odin.accounting.domain.models import Wallet


class WalletCreator(Entity):
    _balance = fields.DecimalField(required=True, private=True)
    _name = fields.StrField(required=True, private=True)

    def __init__(self, wallet_repository: WalletRepository, **kwargs):
        super().__init__(**kwargs)
        self._repository = wallet_repository

    def create(self) -> Wallet:
        wallet = Wallet(
            name=self._name,
            balance=self._balance
        )
        self._repository.add(wallet)
        return wallet
