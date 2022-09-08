import subprocess

import edgedb
import pytest


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
