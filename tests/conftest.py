import uvloop
from pytest import fixture
from starlette.testclient import TestClient

from odin.main import app


@fixture
def test_client():
    return TestClient(app=app)


@fixture(scope='session')
def event_loop():
    loop = uvloop.new_event_loop()
    yield loop
    loop.close()
