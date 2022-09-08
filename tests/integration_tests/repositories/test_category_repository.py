from odin.accounting.models import Category
from odin.accounting.repositories.edgedb_repositories import EdgeDBCategoryRepository


def test_get_by_name(db_client):
    category = Category(name='test')
    repository = EdgeDBCategoryRepository()
    repository.add(category)
    fetched_category = repository.get_by_name(name=category.name)

    assert category.name == fetched_category.name
