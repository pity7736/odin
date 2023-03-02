from odin import settings
from odin.accounting.repositories import InMemoryCategoryRepository
from odin.accounting.repositories.edgedb_repositories import EdgeDBCategoryRepository


def get_category_repository():
    if settings.REPOSITORY == 'in-memory':
        return InMemoryCategoryRepository()
    elif settings.REPOSITORY == 'edgedb':
        return EdgeDBCategoryRepository()
