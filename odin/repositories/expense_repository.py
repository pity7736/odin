from odin.models import Expense
from .exceptions import DoesNotExist

from .category_repository import CategoryRepository


class ExpenseRepository:

    _expenses: dict[str, dict[str, str]] = {}

    def __init__(self):
        self._category_repository = CategoryRepository()

    def add(self, expense: Expense):
        self.__class__._expenses[expense.uuid] = {
            'uuid': expense.uuid,
            'amount': expense.amount,
            'date': expense.date,
            'category_name': expense.category.name
        }
        self._add_category(expense)

    def _add_category(self, expense):
        try:
            self._category_repository.add(expense.category)
        except ValueError:
            pass

    def get_by(self, uuid) -> Expense:
        try:
            expense_data = self._expenses[uuid].copy()
            category = self._category_repository.get_by_name(expense_data.pop('category_name'))
            expense_data['category'] = category
            return Expense(**expense_data)
        except KeyError:
            raise DoesNotExist('Expense not found')

    def get_all(self) -> tuple[Expense]:
        expenses = []
        for expense_data in self._expenses.values():
            data = expense_data.copy()
            data['category'] = self._category_repository.get_by_name(data.pop('category_name'))
            expenses.append(Expense(**data))
        return tuple(expenses)
