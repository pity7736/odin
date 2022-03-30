
class ExpenseRepository:

    _expenses = {}

    def add(self, expense):
        self._expenses[expense.uuid] = expense

    def get_by(self, uuid):
        return self._expenses[uuid]

    def get_all(self):
        pass
