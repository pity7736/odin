from odin.accounting.models import Transfer
from ..repositories import TransferRepository


class InMemoryTransferRepository(TransferRepository):

    _transfers: dict[str, Transfer] = {}

    def add(self, transference):
        self.__class__._transfers[transference.uuid] = transference

    def get_all(self):
        return tuple(self._transfers.values())

    def get_by_uuid(self, uuid):
        return self._transfers.get(uuid)
