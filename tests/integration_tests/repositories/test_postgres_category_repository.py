from pytest import mark

from odin.accounting.application.use_cases import CategoryCreator
from odin.accounting.domain import CategoryType
from odin.accounting.infrastructure.repositories.postgres_repositories import PostgresCategoryRepository


@mark.asyncio
async def test_get_by_name_and_user(db_connection, user_fixture):
    repository = PostgresCategoryRepository()
    category = await CategoryCreator(
        name='test',
        user=user_fixture,
        type=CategoryType.EXPENSE,
        category_repository=repository
    ).create()

    fetched_category = await repository.get_by_name_and_user(category.name, category.user)

    assert fetched_category.id == category.id
