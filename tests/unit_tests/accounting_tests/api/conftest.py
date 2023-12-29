from pytest import fixture

from tests.factories import WalletBuilder


@fixture
def wallet(wallet_repository):
    return WalletBuilder().create()
