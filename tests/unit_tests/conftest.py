from _pytest.fixtures import fixture

from odin.accounting.repositories import ExpenseRepository, WalletRepository, CategoryRepository, TransferenceRepository


@fixture
def db_transaction():
    yield
    ExpenseRepository._expenses.clear()
    WalletRepository._wallets.clear()
    CategoryRepository._categories.clear()
    TransferenceRepository._transfers.clear()
