from odin.accounting.repositories.edgedb_repositories import EdgeDBWalletRepository
from tests.factories import WalletBuilder


def test_get_by_name(db_client):
    wallet = WalletBuilder().build()
    repository = EdgeDBWalletRepository()
    repository.add(wallet)

    fetched_wallet = repository.get_by_name(wallet.name)

    assert fetched_wallet.name == wallet.name
