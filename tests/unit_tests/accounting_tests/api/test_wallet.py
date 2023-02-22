import re

from tests.utils import UUID_PATTERN


def test_create_wallet(test_client, token_value_fixture):
    response = test_client.post(
        '/accounting/wallets',
        json={
            'name': 'test wallet',
            'balance': '10000000',
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.headers['content-type'] == 'application/json'
    assert response.status_code == 201
    assert response_data['name'] == 'test wallet'
    assert response_data['balance'] == '10000000'
    assert re.match(UUID_PATTERN, response_data['uuid'])


def test_create_wallet_with_existing_name(test_client, token_value_fixture):
    test_client.post(
        '/accounting/wallets',
        json={
            'name': 'test wallet',
            'balance': '10000000',
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response = test_client.post(
        '/accounting/wallets',
        json={
            'name': 'test wallet',
            'balance': '10000000',
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )

    assert response.status_code == 400
    assert response.headers['content-type'] == 'application/json'


def test_get_wallet(test_client, token_value_fixture):
    post_response = test_client.post(
        '/accounting/wallets',
        json={
            'name': 'test wallet',
            'balance': '10000000'
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    post_response_data = post_response.json()
    response = test_client.get(
        f'/accounting/wallets/{post_response_data["name"]}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['name'] == 'test wallet'
    assert response_data['balance'] == '10000000'
