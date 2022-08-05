from odin.models import Wallet


class WalletRepository:

    _wallets: dict[str, dict] = {}

    def add(self, wallet: Wallet):
        self.__class__._wallets[wallet.name] = {
            'name': wallet.name,
            'balance': wallet.balance,
            'uuid': wallet.uuid
        }

    def get_by_name(self, name: str) -> Wallet:
        wallet_data = self._wallets.get(name)
        if wallet_data:
            return Wallet(**wallet_data)
