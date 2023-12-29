import uuid

from nyoibo import Entity, fields

from odin.accounting.domain.models import Income, Category, Wallet
from ..repositories import WalletRepository


class IncomeCreator(Entity):
    _date = fields.StrField(private=True)
    _amount = fields.DecimalField(private=True)
    _category = fields.LinkField(to=Category, required=True)
    _wallet: Wallet = fields.LinkField(to=Wallet)

    def __init__(self, wallet_repository: WalletRepository, **kwargs):
        if kwargs.get('category') is None:
            raise ValueError('category is required')
        super().__init__(**kwargs)
        self._repository = wallet_repository

    def create(self) -> Income:
        income = Income(
            date=self._date,
            amount=self._amount,
            category=self._category,
            id=uuid.uuid4()
        )
        self._add_income_to_wallet(income)
        return income

    def _add_income_to_wallet(self, income):
        self._wallet.add_income(income)
        self._repository.add_income(wallet=self._wallet, income=income)
