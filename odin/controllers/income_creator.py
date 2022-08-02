from nyoibo import Entity, fields

from odin.models import Income, Category
from odin.repositories.income_repository import IncomeRepository


class IncomeCreator(Entity):
    _date = fields.StrField(private=True)
    _amount = fields.DecimalField(private=True)
    _category = fields.LinkField(to=Category, required=True)

    def __init__(self, **kwargs):
        if kwargs.get('category') is None:
            raise ValueError('category is required')
        super().__init__(**kwargs)

    def create(self) -> Income:
        income = Income(
            date=self._date,
            amount=self._amount,
            category=self._category
        )
        repository = IncomeRepository()
        repository.add(income)
        return income
