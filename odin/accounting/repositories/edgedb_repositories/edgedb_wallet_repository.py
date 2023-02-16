from typing import Optional

from odin.accounting.models import Wallet

from .db_client import DBClient


class EdgeDBWalletRepository:

    def __init__(self):
        self._client = DBClient()

    def add(self, wallet: Wallet):
        self._client.execute(
            'insert Wallet {name := <str>$name, balance := <decimal>$balance}',
            name=wallet.name,
            balance=wallet.balance
        )

    def get_by_name(self, name) -> Optional[Wallet]:
        record = self._client.query_single('select Wallet {id, name, balance} filter .name = <str>$name', name=name)
        if record:
            return Wallet(
                name=record.name,
                balance=record.balance,
                uuid=record.id
            )
