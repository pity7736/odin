from odin import settings


if settings.REPOSITORY == 'in-memory':
    from .in_memory_reposotiries import InMemoryCategoryRepository, InMemoryWalletRepository,\
        InMemoryExpenseRepository, InMemoryTransferenceRepository
    CategoryRepository = InMemoryCategoryRepository
    ExpenseRepository = InMemoryExpenseRepository
    WalletRepository = InMemoryWalletRepository
    TransferenceRepository = InMemoryTransferenceRepository
