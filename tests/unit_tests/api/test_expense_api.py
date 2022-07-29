import datetime
import re

from pytest import mark

from odin.controllers import WalletCreator
from tests.factories import ExpenseFactory
from tests.utils import UUID_PATTERN


# TODO: test error messages text


def test_create_expense(test_client, category_fixture):
    wallet = WalletCreator(balance='1_000_000', name='savings account').create()
    response = test_client.post(
        '/expenses',
        json={
            'date': '2022-03-27',
            'amount': '100000',
            'category': category_fixture.name,
            'wallet': wallet.name
        }
    )
    response_data = response.json()

    assert response.headers['content-type'] == 'application/json'
    assert response.status_code == 201
    assert response_data['date'] == '2022-03-27'
    assert response_data['amount'] == '100000'
    assert response_data['category'] == category_fixture.name
    assert re.match(UUID_PATTERN, response_data['uuid'])


def test_create_expense_with_wrong_category_name(test_client, category_fixture):
    response = test_client.post(
        '/expenses',
        json={
            'date': '2022-03-27',
            'amount': '100000',
            'category': 'wrong category name'
        }
    )

    assert response.headers['content-type'] == 'application/json'
    assert response.status_code == 400


create_expense_data_params = (
    {'date': '2022-03-27'},
    {'amount': '1000'},
    {'date': 'wrong date', 'amount': '100'},
    {'date': '2022-03-29', 'amount': 'wrong amount'},
)


@mark.parametrize('data', create_expense_data_params)
def test_create_expense_with_missing_or_wrong_data(data, test_client, db_transaction):
    response = test_client.post(
        '/expenses',
        json=data
    )
    assert response.status_code == 400


def test_create_expense_with_date_in_the_future(test_client):
    response = test_client.post(
        '/expenses',
        json={
            'date': (datetime.date.today() + datetime.timedelta(days=1)).isoformat(),
            'amount': '100000'
        }
    )

    assert response.status_code == 400


def test_get_expense(test_client, category_fixture):
    wallet = WalletCreator(balance='1_000_000', name='savings account').create()
    post_response = test_client.post(
        '/expenses',
        json={
            'date': '2022-03-27',
            'amount': '100000',
            'category': category_fixture.name,
            'wallet': wallet.name
        }
    )
    post_response_data = post_response.json()
    response = test_client.get(f'/expenses/{post_response_data["uuid"]}')
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['date'] == post_response_data['date']
    assert response_data['amount'] == post_response_data['amount']
    assert response_data['category'] == post_response_data['category']
    assert response_data['uuid'] == post_response_data['uuid']


def test_get_non_existing_expense(expense_fixture, test_client, db_transaction):
    response = test_client.get('/expenses/1234')
    response_data = response.json()

    assert response.status_code == 404
    assert response.headers['content-type'] == 'application/json'
    assert response_data == {}


def test_get_all_expenses(test_client, db_transaction):
    expenses = ExpenseFactory.create_batch(5)
    expected_expenses = [{
        'date': expense.date.isoformat(),
        'amount': str(expense.amount),
        'uuid': expense.uuid,
        'category': expense.category.name
    } for expense in expenses]
    response = test_client.get('/expenses')
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['expenses'] == expected_expenses
