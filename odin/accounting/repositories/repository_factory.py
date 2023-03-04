from odin import settings
from .in_memory_reposotiries import InMemoryCategoryRepository, InMemoryWalletRepository, InMemoryTransferenceRepository
from .edgedb_repositories import EdgeDBCategoryRepository, EdgeDBWalletRepository
from .repositories import CategoryRepository, WalletRepository, TransferenceRepository


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


def get_transference_repository() -> TransferenceRepository:
    return InMemoryTransferenceRepository()
