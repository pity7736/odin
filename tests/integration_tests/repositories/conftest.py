import subprocess

import edgedb
import pytest

from odin.accounting.repositories.edgedb_repositories.db_client import DBClient


class Rollback(Exception):
    pass


@pytest.fixture(scope='session')
def prepare_db():
    subprocess.run('edgedb -d odin_tests migrate'.split())


@pytest.fixture(scope='session')
def db_client(prepare_db):
    return edgedb.create_client(database='odin_tests')


@pytest.fixture
def db_transaction(mocker, db_client):
    DBClient._instance = None
    try:
        for transaction in db_client.transaction():
            with transaction:
                mocker.patch.object(edgedb, 'create_client', return_value=transaction)
                yield transaction
                raise Rollback
    except Rollback:
        pass
