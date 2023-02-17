from odin.accounting.models import Wallet


class InMemoryWalletRepository:

    _wallets: dict[str, dict] = {}

    def update(self, wallet: Wallet):
        if self.get_by_name(wallet.name):
            self.add(wallet)

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
            return Wallet(**wallet_data)
