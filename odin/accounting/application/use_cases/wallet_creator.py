from nyoibo import Entity, fields

from odin.accounts.domain import User
from odin.accounting.domain import Wallet
from ..repositories import WalletRepository


class WalletCreator(Entity):
    _balance = fields.DecimalField(required=True, private=True)
    _name = fields.StrField(required=True, private=True)
    _user = fields.LinkField(to=User, required=True)

    def __init__(self, wallet_repository: WalletRepository, **kwargs):
        super().__init__(**kwargs)
        self._repository = wallet_repository

    def create(self) -> Wallet:
        wallet = Wallet(
            name=self._name,
            balance=self._balance,
            user=self._user
        )
        self._repository.add(wallet)
        return wallet
