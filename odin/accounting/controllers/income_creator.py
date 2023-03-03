from nyoibo import Entity, fields

from odin.accounting.models import Income, Category, Wallet
from odin.accounting.repositories.repository_factory import get_wallet_repository


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
        return income

    def _add_income_to_wallet(self, income):
        self._wallet.add_income(income)
        wallet_repository = get_wallet_repository()
        wallet_repository.add_income(wallet=self._wallet, income=income)
