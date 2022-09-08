from typing import Optional

from odin.accounting.models import Wallet

import edgedb


class DBClient:

    _instance = None

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._client = edgedb.create_client()
        return cls._instance

    def __getattribute__(self, item):
        if item != '_client':
            return getattr(self._client, item)
        return super().__getattribute__(item)


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
