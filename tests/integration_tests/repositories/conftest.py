import subprocess

import edgedb
import pytest


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
    try:
        for transaction in db_client.transaction():
            with transaction:
                mocker.patch.object(edgedb, 'create_client', return_value=transaction)
                yield transaction
                raise Rollback
    except Rollback:
        pass
