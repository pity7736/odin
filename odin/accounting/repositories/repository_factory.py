from odin import settings
from .in_memory_reposotiries import InMemoryCategoryRepository, InMemoryWalletRepository
from .edgedb_repositories import EdgeDBCategoryRepository, EdgeDBWalletRepository


def get_category_repository():
    if settings.REPOSITORY == 'in-memory':
        return InMemoryCategoryRepository()
    elif settings.REPOSITORY == 'edgedb':
        return EdgeDBCategoryRepository()


def get_wallet_repository():
    if settings.REPOSITORY == 'in-memory':
        return InMemoryWalletRepository()
    elif settings.REPOSITORY == 'edgedb':
        return EdgeDBWalletRepository()
