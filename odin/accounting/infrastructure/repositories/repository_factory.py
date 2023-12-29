from odin.accounting.application.repositories import CategoryRepository, WalletRepository, TransferRepository
from .postgres_repositories import PostgresCategoryRepository, PostgresWalletRepository, PostgresTransferRepository


def get_category_repository() -> CategoryRepository:
    return PostgresCategoryRepository()


def get_wallet_repository() -> WalletRepository:
    return PostgresWalletRepository()


def get_transfer_repository() -> TransferRepository:
    return PostgresTransferRepository()
