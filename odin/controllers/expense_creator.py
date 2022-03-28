import uuid

from nyoibo import Entity, fields

from odin.models import Expense
from odin.repositories import ExpenseRepository


class ExpenseCreator(Entity):
    _date = fields.DateField(private=True)
    _amount = fields.DecimalField(private=True)

    def create(self):
        expense = Expense(
            uuid=str(uuid.uuid4()),
            date=self._date,
            amount=self._amount
        )
        repository = ExpenseRepository()
        repository.add(expense)
        return expense
