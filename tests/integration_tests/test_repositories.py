import subprocess

import edgedb
import pytest

from odin.accounting.models import Category
from odin.accounting.repositories.edgedb_repositories import EdgeDBCategoryRepository


class Rollback(Exception):
    pass


@pytest.fixture(scope='session')
def prepare_db():
    subprocess.run('edgedb migrate'.split())


@pytest.fixture
def db_client(mocker, prepare_db):
    client = edgedb.create_client(database='odin_tests')
    try:
        for transaction in client.transaction():
            with transaction:
                mocker.patch.object(edgedb, 'create_client', return_value=transaction)
                yield transaction
                raise Rollback
    except Rollback:
        pass


def test_get_by_name(db_client):
    category = Category(name='test')
    repository = EdgeDBCategoryRepository()
    repository.add(category)
    fetched_category = repository.get_by_name(name=category.name)

    assert category.name == fetched_category.name
