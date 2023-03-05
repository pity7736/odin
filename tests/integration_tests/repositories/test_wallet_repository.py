from odin.accounting.repositories.edgedb_repositories import EdgeDBWalletRepository
from tests.factories import WalletBuilder


def test_get_wallet_with_expenses_and_with_previous_incomes(db_transaction):
    WalletBuilder().name('another wallet').add_expense('1000').create()
    wallet = WalletBuilder() \
        .add_income('100000') \
        .add_income('100000') \
        .add_expense('50000') \
        .add_expense('20000') \
        .create()
    repository = EdgeDBWalletRepository()
    wallet = repository.get_by_name_with_expenses(wallet.name)

    assert len(wallet.expenses) == 2


def test_get_wallet_with_incomes_and_with_previous_expenses(db_transaction):
    WalletBuilder().name('another wallet').add_income('1000').create()
    wallet = WalletBuilder() \
        .add_income('100000') \
        .add_income('100000') \
        .add_expense('50000') \
        .add_expense('20000') \
        .create()
    repository = EdgeDBWalletRepository()
    wallet = repository.get_by_name_with_incomes(wallet.name)

    assert len(wallet.incomes) == 2
