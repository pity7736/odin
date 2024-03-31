from pytest import raises, mark

from odin.accounting.application.use_cases import CategoryCreator
from odin.accounting.domain import CategoryType
from tests.factories import UserFactory, CategoryFactory


@mark.asyncio
async def test_create_category_with_same_name_and_user(category_repository):
    user = UserFactory.create()
    category_name = 'test'
    await CategoryFactory.create(user=user, name=category_name)
    category_creator = CategoryCreator(
        name=category_name,
        user=user,
        type=CategoryType.EXPENSE,
        category_repository=category_repository
    )
    with raises(ValueError) as error:
        await category_creator.create()

    assert str(error.value) == f'there is already a category with name {category_name}'


@mark.asyncio
async def test_create_category_with_same_name_and_different_user(category_repository):
    user = UserFactory.create()
    await CategoryFactory.create(user=user, name='test')
    category_repository.get_by_name_and_user.return_value = None
    category_creator = CategoryCreator(
        name='test',
        user=UserFactory.create(),
        type=CategoryType.EXPENSE,
        category_repository=category_repository
    )
    category = await category_creator.create()

    assert category.name == 'test'
