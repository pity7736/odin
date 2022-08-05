from odin.models import Wallet

from .expense_repository import ExpenseRepository


class WalletRepository:

    _wallets: dict[str, dict] = {}

    def add(self, wallet: Wallet):
        self.__class__._wallets[wallet.name] = {
            'name': wallet.name,
            'balance': wallet.balance,
            'uuid': wallet.uuid,
            'expenses_uuid': [expense.uuid for expense in wallet.expenses]
        }

    def get_by_name(self, name: str) -> Wallet:
        wallet_data = self._wallets.get(name)
        if wallet_data:
            expenses = []
            for expense_uuid in wallet_data['expenses_uuid']:
                expenses.append(ExpenseRepository().get_by(uuid=expense_uuid))
            wallet_data['expenses'] = expenses
            return Wallet(**wallet_data)
