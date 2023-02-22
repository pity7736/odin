from pytest import fixture

from tests.factories import WalletBuilder


@fixture
def wallet():
    return WalletBuilder().create()
