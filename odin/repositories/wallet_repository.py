from odin.models import Wallet


class WalletRepository:

    _wallets = {}

    def add(self, wallet: Wallet):
        self.__class__._wallets[wallet.name] = wallet

    def get_by_name(self, name: str) -> Wallet:
        return self._wallets.get(name)
