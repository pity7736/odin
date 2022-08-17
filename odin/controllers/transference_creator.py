import datetime

from nyoibo import Entity, fields
from nyoibo.fields import Decimal

from odin.models import Wallet, Transference
from odin.repositories import CategoryRepository
from .expense_creator import ExpenseCreator
from .income_creator import IncomeCreator


class TransferenceCreator(Entity):
    _source = fields.LinkField(to=Wallet, private=True)
    _target = fields.LinkField(to=Wallet, private=True)

    def transfer(self, amount: Decimal, date: datetime.date = None):
        date = date or datetime.date.today()
        transference_category = CategoryRepository().get_by_name('transference')
        expense = self._create_expense(amount, date, transference_category)
        income = self._create_income(amount, date, transference_category)
        return Transference(
            source=self._source,
            target=self._target,
            expense=expense,
            income=income,
            amount=amount,
            date=date
        )

    def _create_expense(self, amount, date, transference_category):
        return ExpenseCreator(
            amount=amount,
            date=date,
            category=transference_category,
            wallet=self._source
        ).create()

    def _create_income(self, amount, date, transference_category):
        return IncomeCreator(
            amount=amount,
            date=date,
            category=transference_category,
            wallet=self._target
        ).create()
