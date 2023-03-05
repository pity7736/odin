from odin import settings
from .in_memory_reposotiries import InMemoryCategoryRepository, InMemoryWalletRepository, InMemoryTransferRepository
from .edgedb_repositories import EdgeDBCategoryRepository, EdgeDBWalletRepository, EdgeDBTransferRepository
from .repositories import CategoryRepository, WalletRepository, TransferRepository


def get_category_repository() -> CategoryRepository:
    if settings.REPOSITORY == 'in-memory':
        return InMemoryCategoryRepository()
    elif settings.REPOSITORY == 'edgedb':
        return EdgeDBCategoryRepository()


def get_wallet_repository() -> WalletRepository:
    if settings.REPOSITORY == 'in-memory':
        return InMemoryWalletRepository()
    elif settings.REPOSITORY == 'edgedb':
        return EdgeDBWalletRepository()


def get_transfer_repository() -> TransferRepository:
    if settings.REPOSITORY == 'in-memory':
        return InMemoryTransferRepository()
    elif settings.REPOSITORY == 'edgedb':
        return EdgeDBTransferRepository()
