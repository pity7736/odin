import datetime
import uuid

from nyoibo import Entity, fields
from nyoibo.fields import Decimal

from odin.accounting.domain.models import Wallet, Transfer
from .expense_creator import ExpenseCreator
from .income_creator import IncomeCreator
from ..repositories import WalletRepository, TransferRepository, CategoryRepository


class TransferCreator(Entity):
    _source = fields.LinkField(to=Wallet, private=True)
    _target = fields.LinkField(to=Wallet, private=True)

    def __init__(self, wallet_repository: WalletRepository, transfer_repository: TransferRepository,
                 category_repository: CategoryRepository, **kwargs):
        if kwargs.get('source') is None:
            raise ValueError('source is required')

        if kwargs.get('target') is None:
            raise ValueError('target is required')
        super().__init__(**kwargs)
        self._wallet_repository = wallet_repository
        self._transfer_repository = transfer_repository
        self._category_repository = category_repository

    @classmethod
    async def from_wallet_names(cls, source_name: str, target_name: str, wallet_repository: WalletRepository,
                                transfer_repository: TransferRepository, category_repository: CategoryRepository):
        return cls(
            source=await wallet_repository.get_by_name(source_name),
            target=await wallet_repository.get_by_name(target_name),
            wallet_repository=wallet_repository,
            transfer_repository=transfer_repository,
            category_repository=category_repository
        )

    async def transfer(self, amount: Decimal, date: datetime.date = None):
        return await self._create_transfer(amount, date or datetime.date.today())

    async def _create_transfer(self, amount: Decimal, date: datetime.date):
        category = await self._category_repository.get_by_name_and_user('transfer', self._source.user)
        transfer = Transfer(
            source=self._source,
            target=self._target,
            expense=await self._create_expense(amount, date, category),
            income=await self._create_income(amount, date, category),
            amount=amount,
            date=date,
            id=uuid.uuid4()
        )
        await self._transfer_repository.add(transfer)
        return transfer

    async def _create_expense(self, amount, date, category):
        return await ExpenseCreator(
            amount=amount,
            date=date,
            category=category,
            wallet=self._source,
            wallet_repository=self._wallet_repository
        ).create()

    async def _create_income(self, amount, date, category):
        return await IncomeCreator(
            amount=amount,
            date=date,
            category=category,
            wallet=self._target,
            wallet_repository=self._wallet_repository
        ).create()
