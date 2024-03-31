from pytest_asyncio import fixture as async_fixture

from tests.factories import WalletBuilder


@async_fixture
async def wallet_fixture(wallet_repository, category_repository):
    return await WalletBuilder().create()
