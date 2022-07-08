from odin.models import Expense
from odin.repositories import ExpenseRepository
from odin.repositories.exceptions import DoesNotExist


class ExpenseGetter:

    __slots__ = ('_repository',)

    def __init__(self):
        self._repository = ExpenseRepository()

    def get_by_uuid(self, uuid: str) -> Expense | None:
        try:
            return self._repository.get_by(uuid=uuid)
        except DoesNotExist:
            return None

    def all(self) -> tuple[Expense]:
        return self._repository.get_all()
