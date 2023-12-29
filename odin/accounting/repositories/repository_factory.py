from .in_memory_reposotiries import InMemoryCategoryRepository, InMemoryWalletRepository, InMemoryTransferRepository
from .repositories import CategoryRepository, WalletRepository, TransferRepository


def get_category_repository() -> CategoryRepository:
    return InMemoryCategoryRepository()


def get_wallet_repository() -> WalletRepository:
    return InMemoryWalletRepository()


def get_transfer_repository() -> TransferRepository:
    return InMemoryTransferRepository()
