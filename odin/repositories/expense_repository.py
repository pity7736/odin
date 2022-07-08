from odin.models import Expense
from .exceptions import DoesNotExist


class ExpenseRepository:

    _expenses: dict[str, Expense] = {}

    def add(self, expense: Expense):
        self._expenses[expense.uuid] = expense

    def get_by(self, uuid) -> Expense:
        try:
            return self._expenses[uuid]
        except KeyError:
            raise DoesNotExist('Expense not found')

    def get_all(self) -> tuple[Expense]:
        return tuple(self._expenses.values())
