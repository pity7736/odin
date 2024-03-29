from decimal import Decimal

from pytest import mark

from odin.accounting.application.use_cases import TransferCreator
from odin.accounting.infrastructure.repositories.postgres_repositories import PostgresTransferRepository, \
    PostgresCategoryRepository
from tests.factories import WalletBuilder, CategoryFactory


@mark.asyncio
async def test_get_by_id(db_connection, user_fixture, wallet_repository):
    await CategoryFactory.create(name='transfer', user=user_fixture)
    wallet_a = await WalletBuilder().user(user_fixture).create()
    wallet_b = await WalletBuilder().user(user_fixture).create()
    repository = PostgresTransferRepository()
    transfer = await TransferCreator(
        wallet_repository=wallet_repository,
        transfer_repository=repository,
        category_repository=PostgresCategoryRepository(),
        source=wallet_a,
        target=wallet_b
    ).transfer(Decimal('100_000'))

    fetched_transfer = await repository.get_by_id(transfer.id)

    assert fetched_transfer == transfer
    assert fetched_transfer.source == wallet_a
    assert fetched_transfer.target == wallet_b
