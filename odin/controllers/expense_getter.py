from odin.repositories import ExpenseRepository


class ExpenseGetter:

    __slots__ = ('_repository',)

    def __init__(self):
        self._repository = ExpenseRepository()

    def get_by_uuid(self, uuid):
        try:
            return self._repository.get_by(uuid=uuid)
        except KeyError:
            return None

    def all(self):
        return self._repository.get_all()
