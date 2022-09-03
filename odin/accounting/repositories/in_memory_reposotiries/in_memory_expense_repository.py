from odin.accounting.models import Expense
from ..exceptions import DoesNotExist

from .in_memory_category_repository import InMemoryCategoryRepository


class InMemoryExpenseRepository:

    _expenses: dict[str, dict[str, str]] = {}

    def __init__(self):
        self._category_repository = InMemoryCategoryRepository()

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
            expense_data = self._expenses[uuid]
        except KeyError:
            raise DoesNotExist('Expense not found')
        else:
            return Expense(
                **expense_data,
                category=self._category_repository.get_by_name(expense_data.get('category_name'))
            )

    def get_all(self) -> tuple[Expense]:
        expenses = []
        for expense_data in self._expenses.values():
            expenses.append(
                Expense(
                    **expense_data,
                    category=self._category_repository.get_by_name(expense_data.get('category_name')))
            )
        return tuple(expenses)
