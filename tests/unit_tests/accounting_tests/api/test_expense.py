import datetime
import re

from pytest import mark

from tests.factories import ExpenseFactory, WalletBuilder
from tests.utils import UUID_PATTERN


# TODO: test error messages text


def test_create_expense(test_client, category_fixture, wallet, token_value_fixture):
    response = test_client.post(
        f'/accounting/wallets/{wallet.name}/expenses',
        json={
            'date': '2022-03-27',
            'amount': '100000',
            'category': category_fixture.name,
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.headers['content-type'] == 'application/json'
    assert response.status_code == 201
    assert response_data['date'] == '2022-03-27'
    assert response_data['amount'] == '100000'
    assert response_data['category'] == category_fixture.name
    assert re.match(UUID_PATTERN, response_data['uuid'])


def test_create_expense_with_wrong_category_name(test_client, category_fixture, wallet, token_value_fixture):
    response = test_client.post(
        f'/accounting/wallets/{wallet.name}/expenses',
        json={
            'date': '2022-03-27',
            'amount': '100000',
            'category': 'wrong category name'
        },
        headers={'Authorization': f'token {token_value_fixture}'}
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
def test_create_expense_with_missing_or_wrong_data(data, test_client, wallet, token_value_fixture):
    response = test_client.post(
        f'/accounting/wallets/{wallet.name}/expenses',
        json=data,
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    assert response.status_code == 400


def test_create_expense_with_date_in_the_future(test_client, wallet, token_value_fixture):
    response = test_client.post(
        f'/accounting/wallets/{wallet.name}/expenses',
        json={
            'date': (datetime.date.today() + datetime.timedelta(days=1)).isoformat(),
            'amount': '100000'
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )

    assert response.status_code == 400


def test_create_expense_with_higher_amount_that_wallet_balance(test_client, category_fixture, wallet,
                                                               token_value_fixture):
    wallet = WalletBuilder().balance('100_000').create()
    response = test_client.post(
        f'/accounting/wallets/{wallet.name}/expenses',
        json={
            'date': '2022-03-27',
            'amount': '102000',
            'category': category_fixture.name,
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 400
    assert response_data['error'] == 'expense amount must be lower than wallet balance'


def test_get_expense(test_client, category_fixture, wallet, token_value_fixture):
    post_response = test_client.post(
        f'/accounting/wallets/{wallet.name}/expenses',
        json={
            'date': '2022-03-27',
            'amount': '100000',
            'category': category_fixture.name,
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    post_response_data = post_response.json()
    response = test_client.get(
        f'/accounting/wallets/{wallet.name}/expenses/{post_response_data["uuid"]}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['date'] == post_response_data['date']
    assert response_data['amount'] == post_response_data['amount']
    assert response_data['category'] == post_response_data['category']
    assert response_data['uuid'] == post_response_data['uuid']


def test_get_non_existing_expense(expense_fixture, test_client, wallet, token_value_fixture):
    response = test_client.get(
        f'/accounting/wallets/{wallet.name}/expenses/1234',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 404
    assert response.headers['content-type'] == 'application/json'
    assert response_data == {}


def test_get_all_expenses(test_client, token_value_fixture, wallet):
    expenses = ExpenseFactory.create_batch(5)
    expected_expenses = [{
        'date': expense.date.isoformat(),
        'amount': str(expense.amount),
        'uuid': expense.uuid,
        'category': expense.category.name
    } for expense in expenses]
    response = test_client.get(
        f'/accounting/wallets/{wallet.name}/expenses',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['expenses'] == expected_expenses
