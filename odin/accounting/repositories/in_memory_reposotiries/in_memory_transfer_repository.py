from uuid import uuid4

from odin.accounting.models import Transfer
from ..repositories import TransferRepository


class InMemoryTransferRepository(TransferRepository):

    _transfers: dict[str, Transfer] = {}

    def add(self, transfer):
        transfer.uuid = uuid4()
        self.__class__._transfers[transfer.uuid] = transfer

    def get_all(self):
        return tuple(self._transfers.values())

    def get_by_uuid(self, uuid):
        return self._transfers.get(uuid)
