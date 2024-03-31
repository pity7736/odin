import uvloop
from pytest import fixture
from starlette.testclient import TestClient

from odin import settings
from odin.main import app


@fixture
def test_client():
    return TestClient(app=app)


@fixture(scope='session')
def event_loop():
    loop = uvloop.new_event_loop()
    yield loop
    loop.close()


@fixture(autouse=True, scope='session')
def db_name_setting():
    db_name = settings.DB_NAME
    settings.DB_NAME = f'{db_name}_tests'
    yield db_name
    settings.DB_NAME = db_name
