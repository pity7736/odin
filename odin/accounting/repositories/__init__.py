from odin import settings


if settings.REPOSITORY == 'in-memory':
    from .in_memory_reposotiries import InMemoryCategoryRepository, InMemoryWalletRepository,\
        InMemoryTransferenceRepository
    CategoryRepository = InMemoryCategoryRepository
    WalletRepository = InMemoryWalletRepository
    TransferenceRepository = InMemoryTransferenceRepository
