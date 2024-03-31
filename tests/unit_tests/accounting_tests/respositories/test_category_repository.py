# from pytest import mark, fixture, raises
#
# from odin.accounting.models import Category
# from odin.accounting.repositories.in_memory_reposotiries import InMemoryCategoryRepository
# from tests.factories import CategoryFactory
#
#
# @fixture
# def repository_fixture(db_transaction):
#     return InMemoryCategoryRepository()
#
#
# def create_category(category_name):
#     return CategoryFactory.create(name=category_name)
#
#
# def test_add_category(repository_fixture):
#     category_name = 'test'
#     create_category(category_name)
#
#     category = repository_fixture.get_by_name(category_name)
#
#     assert category.name == category_name
#
#
# def test_add_existing_category(repository_fixture):
#     category_name = 'test'
#     create_category(category_name)
#
#     with raises(ValueError) as error:
#         repository_fixture.add(Category(name=category_name))
#
#     assert str(error.value) == f'a category with name {category_name} already exists'
#
#
# incorrect_category_params = (
#     None,
#     '',
#     'wrong type'
# )
#
#
# @mark.parametrize('category', incorrect_category_params)
# def test_add_incorrect_category(category, repository_fixture):
#     with raises(AssertionError) as error:
#         repository_fixture.add(category)
#
#     assert str(error.value) == 'category argument must be Category instance'
#
#
# category_name_params = (
#     ('Test', 'test'),
#     ('Test', 'Test'),
#     ('Test', 'TeSt'),
#     ('test', 'Test')
# )
#
#
# @mark.parametrize('category_name, lookup_name', category_name_params)
# def test_get_by_name(category_name, lookup_name, repository_fixture):
#     create_category(category_name)
#
#     category = repository_fixture.get_by_name(lookup_name)
#
#     assert category.name == category_name.lower()
#
#
# def test_get_by_wrong_name(category_fixture, repository_fixture):
#     assert repository_fixture.get_by_name('wrong name') is None
#
#
# def test_by_none_name(category_fixture, repository_fixture):
#     assert repository_fixture.get_by_name(None) is None
#
#
# def test_get_all(repository_fixture):
#     CategoryFactory.create_batch(5)
#     categories = repository_fixture.get_all()
#
#     assert len(categories) == 5
#     for category in categories:
#         assert isinstance(category, Category)
