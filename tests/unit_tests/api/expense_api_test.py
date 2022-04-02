import re

from pytest import mark
from starlette.testclient import TestClient

from odin.api import app
from tests.utils import UUID_PATTERN


def test_create_expense(test_client):
    response = test_client.post(
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
def test_create_expense_with_missing_or_wrong_data(data, test_client):
    response = test_client.post(
        '/expenses',
        json=data
    )
    assert response.status_code == 400


def test_get_expense():
    test_client = TestClient(app)
    post_response = test_client.post(
        '/expenses',
        json={
            'date': '2022-03-27',
            'amount': '100000'
        }
    )
    response_data = post_response.json()
    response = test_client.get(f'/expenses/{response_data["uuid"]}')
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['date'] == '2022-03-27'
    assert response_data['amount'] == '100000'


def test_get_non_existing_expense(expense_fixture, test_client):
    response = test_client.get('/expenses/1234')
    response_data = response.json()

    assert response.status_code == 404
    assert response.headers['content-type'] == 'application/json'
    assert response_data == {}
