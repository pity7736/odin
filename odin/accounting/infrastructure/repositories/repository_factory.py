from odin.accounting.application.repositories import CategoryRepository, WalletRepository, TransferRepository
from .postgres_repositories import PostgresCategoryRepository, PostgresWalletRepository, PostgresTransferRepository


class RepositoryFactory:

    _repositories = {
        'wallet': PostgresWalletRepository,
        'category': PostgresCategoryRepository,
        'transfer': PostgresTransferRepository
    }

    def get_wallet_repository(self) -> WalletRepository:
        return self._repositories['wallet']()

    def get_category_repository(self) -> CategoryRepository:
        return self._repositories['category']()

    def get_transfer_repository(self) -> TransferRepository:
        return self._repositories['transfer']()
