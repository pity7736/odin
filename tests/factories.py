import datetime
from decimal import Decimal
from uuid import uuid4

import factory

from odin.models import Expense, Category
from odin.repositories import ExpenseRepository, CategoryRepository


class CategoryFactory(factory.Factory):
    name = factory.Sequence(lambda n: f'test category{n}')

    class Meta:
        model = Category

    @classmethod
    def _create(cls, model_class, *args, **kwargs):
        category = super()._create(model_class, *args, **kwargs)
        repository = CategoryRepository()
        repository.add(category)
        return category


class ExpenseFactory(factory.Factory):
    uuid = factory.LazyFunction(uuid4)
    date = datetime.date(2022, 3, 30)
    amount = Decimal('100_000')
    category = factory.SubFactory(CategoryFactory)

    class Meta:
        model = Expense

    @classmethod
    def _create(cls, model_class, *args, **kwargs):
        expense = super()._create(model_class, *args, **kwargs)
        repository = ExpenseRepository()
        repository.add(expense)
        return expense
