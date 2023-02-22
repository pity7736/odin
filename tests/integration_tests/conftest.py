import subprocess

import edgedb
import pytest
from pytest import fixture

from odin import settings
from odin.accounting.repositories.edgedb_repositories.db_client import DBClient
from odin.accounts.models import User
from odin.accounts.repositories.edgedb_repositories import EdgeDBUserRepository


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
    repository = settings.REPOSITORY
    settings.REPOSITORY = 'edgedb'
    try:
        for transaction in db_client.transaction():
            with transaction:
                mocker.patch.object(edgedb, 'create_client', return_value=transaction)
                yield transaction
                raise Rollback
    except Rollback:
        pass
    settings.REPOSITORY = repository


@fixture
def user_fixture(db_transaction):
    user = User(
        email='me@raiseexception.com',
        password='test',
        first_name='julián',
        last_name='cortés'
    )
    EdgeDBUserRepository().add(user)
    return user


@fixture
def token_value_fixture(user_fixture, test_client):
    login_response = test_client.post(
        '/auth/login',
        json={
            'email': user_fixture.email,
            'password': 'test'
        }
    )
    data = login_response.json()
    return data['token']
