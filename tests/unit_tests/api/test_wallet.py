import re

from tests.utils import UUID_PATTERN


def test_create_wallet(test_client, db_transaction):
    response = test_client.post(
        '/wallets',
        json={
            'name': 'test wallet',
            'balance': '10000000',
        }
    )
    response_data = response.json()

    assert response.headers['content-type'] == 'application/json'
    assert response.status_code == 201
    assert response_data['name'] == 'test wallet'
    assert response_data['balance'] == '10000000'
    assert re.match(UUID_PATTERN, response_data['uuid'])


def test_create_wallet_with_existing_name(test_client, db_transaction):
    test_client.post(
        '/wallets',
        json={
            'name': 'test wallet',
            'balance': '10000000',
        }
    )
    response = test_client.post(
        '/wallets',
        json={
            'name': 'test wallet',
            'balance': '10000000',
        }
    )

    assert response.status_code == 400
    assert response.headers['content-type'] == 'application/json'


def test_get_wallet(test_client, db_transaction):
    post_response = test_client.post(
        '/wallets',
        json={
            'name': 'test wallet',
            'balance': '10000000'
        }
    )
    post_response_data = post_response.json()
    response = test_client.get(f'/wallets/{post_response_data["name"]}')
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['name'] == 'test wallet'
    assert response_data['balance'] == '10000000'
