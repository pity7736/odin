from .exceptions import DoesNotExist


class ExpenseRepository:

    _expenses = {}

    def add(self, expense):
        self._expenses[expense.uuid] = expense

    def get_by(self, uuid):
        try:
            return self._expenses[uuid]
        except KeyError:
            raise DoesNotExist('Expense not found')

    def get_all(self):
        pass
