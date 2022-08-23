from odin.models import Transference


class TransferenceRepository:

    _transfers: dict[str, Transference] = {}

    def add(self, transference: Transference):
        self.__class__._transfers[transference.uuid] = transference

    def get_all(self) -> tuple[Transference]:
        return tuple(self._transfers.values())

    def get_by_uuid(self, uuid):
        return self._transfers.get(uuid)
