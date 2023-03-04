from odin.accounting.models import Transference
from ..repositories import TransferenceRepository


class InMemoryTransferenceRepository(TransferenceRepository):

    _transfers: dict[str, Transference] = {}

    def add(self, transference):
        self.__class__._transfers[transference.uuid] = transference

    def get_all(self):
        return tuple(self._transfers.values())

    def get_by_uuid(self, uuid):
        return self._transfers.get(uuid)
