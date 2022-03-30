import re

from pytest import mark
from starlette.testclient import TestClient

from odin.api import app
from tests.utils import UUID_PATTERN


def test_create_expense():
    client = TestClient(app)
    response = client.post(
        '/expenses',
        json={
            'date': '2022-03-27',
            'amount': '100000'
        }
    )
    response_data = response.json()

    assert response.headers['content-type'] == 'application/json'
    assert response.status_code == 201
    assert response_data['date'] == '2022-03-27'
    assert response_data['amount'] == '100000'
    assert re.match(UUID_PATTERN, response_data['uuid'])


create_expense_data_params = (
    {'date': '2022-03-27'},
    {'amount': '1000'},
    {'date': 'wrong date', 'amount': '100'},
    {'date': '2022-03-29', 'amount': 'wrong amount'},
)


@mark.parametrize('data', create_expense_data_params)
def test_create_expense_with_missing_or_wrong_data(data):
    client = TestClient(app)
    response = client.post(
        '/expenses',
        json=data
    )
    assert response.status_code == 400
