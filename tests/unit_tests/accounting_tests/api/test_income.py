import re

from tests.utils import UUID_PATTERN


def test_create_income(test_client, category_fixture, wallet, token_value_fixture):
    response = test_client.post(
        f'/accounting/wallets/{wallet.name}/incomes',
        json={
            'date': '2022-07-31',
            'amount': '10000000',
            'category': category_fixture.name,
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 201
    assert response.headers['content-type'] == 'application/json'
    assert response_data['amount'] == '10000000'
    assert response_data['date'] == '2022-07-31'
    assert re.match(UUID_PATTERN, response_data['uuid'])


def test_create_income_without_category(test_client, category_fixture, wallet, token_value_fixture):
    response = test_client.post(
        f'/accounting/wallets/{wallet.name}/incomes',
        json={
            'date': '2022-07-31',
            'amount': '10000000',
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )

    assert response.status_code == 400


def test_get_income(test_client, category_fixture, wallet, token_value_fixture):
    post_response = test_client.post(
        f'/accounting/wallets/{wallet.name}/incomes',
        json={
            'date': '2022-07-31',
            'amount': '10000000',
            'category': category_fixture.name,
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    post_response_data = post_response.json()

    response = test_client.get(f'/accounting/wallets/{wallet.name}/incomes/{post_response_data["uuid"]}')
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['amount'] == '10000000'
    assert response_data['date'] == '2022-07-31'
    assert response_data['category'] == category_fixture.name
    assert re.match(UUID_PATTERN, response_data['uuid'])
