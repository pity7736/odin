import uuid
from pathlib import Path

import asyncpg
from pytest import fixture
from pytest_asyncio import fixture as async_fixture

from odin import settings
from odin.accounts.domain import Token, User
from odin.accounts.domain.crypto import get_random_string
from odin.accounts.infrastructure.repositories.postgres_repositories import PostgresTokenRepository, \
    PostgresUserRepository

CURRENT_DIR = Path(__file__).parent


@fixture(scope='session')
def schema():
    with open(CURRENT_DIR / 'db_schema.sql') as f:
        result = f.read()
    return result


@async_fixture(scope='session')
async def create_db():
    db_name = settings.DB_NAME
    connection = await asyncpg.connect(
        host=settings.DB_HOST,
        user=settings.DB_USER,
        database=db_name,
        password=settings.DB_PASSWORD,
        port=settings.DB_PORT,
    )
    settings.DB_NAME = f'{db_name}_tests'
    await connection.execute(f'DROP DATABASE IF EXISTS {settings.DB_NAME}')
    await connection.execute(f'CREATE DATABASE {settings.DB_NAME} WITH OWNER odin')
    yield
    await connection.execute(f'DROP DATABASE {settings.DB_NAME}')
    await connection.close()
    settings.DB_NAME = db_name


@async_fixture(scope='session')
async def db_pool(create_db, schema):
    pool = await asyncpg.create_pool(
        host=settings.DB_HOST,
        user=settings.DB_USER,
        database=settings.DB_NAME,
        password=settings.DB_PASSWORD,
        port=settings.DB_PORT,
        min_size=1
    )
    connection = await pool.acquire()
    await connection.execute(schema)
    await pool.release(connection)
    yield pool
    await pool.close()


@async_fixture
async def db_connection(db_pool):
    connection = await db_pool.acquire()
    yield connection
    await connection.execute('truncate table tokens cascade')
    await db_pool.release(connection)


@fixture
def token_repository():
    return PostgresTokenRepository()


@fixture
def user_repository():
    return PostgresUserRepository()


@async_fixture
async def user_fixture(user_repository):
    user = User(
        email='me@raiseexception.com',
        password='test',
        first_name='julián',
        last_name='cortés',
        id=uuid.uuid4()
    )
    await user_repository.add(user)
    return user


@async_fixture
async def token_value_fixture(user_fixture, test_client, token_repository):
    token = Token(
        value=get_random_string(length=50),
        user=user_fixture
    )
    await token_repository.add(token)
    return token.value
