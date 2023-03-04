import datetime

from nyoibo import Entity, fields
from nyoibo.fields import Decimal

from odin.accounting.models import Wallet, Transfer
from odin.accounting.repositories import CategoryRepository
from .expense_creator import ExpenseCreator
from .income_creator import IncomeCreator
from ..repositories.repository_factory import get_wallet_repository, get_transfer_repository


class TransferCreator(Entity):
    _source = fields.LinkField(to=Wallet, private=True)
    _target = fields.LinkField(to=Wallet, private=True)

    def __init__(self, **kwargs):
        if kwargs.get('source') is None:
            raise ValueError('source is required')

        if kwargs.get('target') is None:
            raise ValueError('target is required')
        super().__init__(**kwargs)

    @classmethod
    def from_wallet_names(cls, source_name: str, target_name: str):
        repository = get_wallet_repository()
        return cls(
            source=repository.get_by_name(source_name),
            target=repository.get_by_name(target_name)
        )

    def transfer(self, amount: Decimal, date: datetime.date = None):
        return self._create_transference(amount, date or datetime.date.today())

    def _create_transference(self, amount: Decimal, date: datetime.date):
        category = CategoryRepository().get_by_name('transference')
        transference = Transfer(
            source=self._source,
            target=self._target,
            expense=self._create_expense(amount, date, category),
            income=self._create_income(amount, date, category),
            amount=amount,
            date=date
        )
        get_transfer_repository().add(transference)
        return transference

    def _create_expense(self, amount, date, category):
        return ExpenseCreator(
            amount=amount,
            date=date,
            category=category,
            wallet=self._source
        ).create()

    def _create_income(self, amount, date, category):
        return IncomeCreator(
            amount=amount,
            date=date,
            category=category,
            wallet=self._target
        ).create()
