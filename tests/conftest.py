from pytest import fixture
from starlette.testclient import TestClient

from odin.main import app
from tests.factories import CategoryFactory


@fixture
def test_client():
    return TestClient(app=app)


@fixture
def transfer_category():
    return CategoryFactory.create(name='transfer')
