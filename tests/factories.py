import datetime
from decimal import Decimal
from uuid import uuid4

import factory

from odin.models import Expense
from odin.repositories import ExpenseRepository


class ExpenseFactory(factory.Factory):
    uuid = factory.LazyFunction(uuid4)
    date = datetime.date(2022, 3, 30)
    amount = Decimal('100_000')

    class Meta:
        model = Expense

    @classmethod
    def _create(cls, model_class, *args, **kwargs):
        expense = super()._create(model_class, *args, **kwargs)
        repository = ExpenseRepository()
        repository.add(expense)
        return expense
