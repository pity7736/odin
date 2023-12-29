from pytest import fixture
from starlette.testclient import TestClient

from odin.main import app


@fixture
def test_client():
    return TestClient(app=app)
