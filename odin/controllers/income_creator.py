from nyoibo import Entity, fields

from odin.models import Income, Category, Wallet
from odin.repositories import WalletRepository
from odin.repositories.income_repository import IncomeRepository


class IncomeCreator(Entity):
    _date = fields.StrField(private=True)
    _amount = fields.DecimalField(private=True)
    _category = fields.LinkField(to=Category, required=True)
    _wallet: Wallet = fields.LinkField(to=Wallet)

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
        self._add_income_to_wallet(income)
        repository = IncomeRepository()
        repository.add(income)
        return income

    def _add_income_to_wallet(self, income):
        self._wallet.add_income(income)
        wallet_repository = WalletRepository()
        wallet_repository.update(wallet=self._wallet)
