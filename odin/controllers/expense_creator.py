import datetime

from nyoibo import Entity, fields

from odin.models import Expense, Category
from odin.repositories import ExpenseRepository


class ExpenseCreator(Entity):
    _date = fields.DateField(private=True, required=True)
    _amount = fields.DecimalField(private=True, required=True)
    _category = fields.LinkField(to=Category, private=True, required=True)

    def __init__(self, **kwargs):
        if kwargs.get('category') is None:
            raise ValueError('category is required')

        super().__init__(**kwargs)
        if self._date > datetime.date.today():
            raise ValueError('date must be less or equal than today.')

    def create(self) -> Expense:
        expense = Expense(
            date=self._date,
            amount=self._amount,
            category=self._category
        )
        repository = ExpenseRepository()
        repository.add(expense)
        return expense
