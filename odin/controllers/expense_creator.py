import datetime
from decimal import Decimal

from nyoibo import Entity, fields

from odin.models import Expense, Category, Wallet
from odin.repositories import ExpenseRepository


class ExpenseCreator(Entity):
    _date: datetime.date = fields.DateField(private=True, required=True)
    _amount: Decimal = fields.DecimalField(private=True, required=True)
    _category: Category = fields.LinkField(to=Category, private=True, required=True)
    _wallet: Wallet = fields.LinkField(to=Wallet, private=True, required=True)

    def __init__(self, **kwargs):
        if kwargs.get('category') is None:
            raise ValueError('category is required')

        if kwargs.get('wallet') is None:
            raise ValueError('wallet is required')

        super().__init__(**kwargs)
        if self._date > datetime.date.today():
            raise ValueError('date must be less or equal than today.')

    def create(self) -> Expense:
        expense = Expense(
            date=self._date,
            amount=self._amount,
            category=self._category,
        )
        try:
            self._wallet.add_expense(expense)
        except AssertionError as error:
            raise ValueError(str(error))
        repository = ExpenseRepository()
        repository.add(expense)
        return expense
