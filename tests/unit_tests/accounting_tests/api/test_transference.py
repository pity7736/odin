import re
import uuid

from tests.factories import WalletBuilder, CategoryFactory
from tests.utils import UUID_PATTERN


def test_create(test_client, token_value_fixture):
    CategoryFactory.create(name='transference')
    wallet_source = WalletBuilder().create()
    wallet_target = WalletBuilder().name('cash').create()
    response = test_client.post(
        '/transfers',
        json={
            'source': wallet_source.name,
            'target': wallet_target.name,
            'amount': '100000',
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 201
    assert response.headers['content-type'] == 'application/json'
    assert response_data['source'] == wallet_source.name
    assert response_data['target'] == wallet_target.name
    assert response_data['amount'] == '100000'
    assert re.match(UUID_PATTERN, response_data['uuid'])


def test_create_with_non_existing_source_wallet(test_client, token_value_fixture):
    CategoryFactory.create(name='transference')
    wallet_target = WalletBuilder().name('cash').create()
    response = test_client.post(
        '/transfers',
        json={
            'source': 'source wallet',
            'target': wallet_target.name,
            'amount': '100000',
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )

    assert response.status_code == 400
    assert response.headers['content-type'] == 'application/json'


def test_create_with_non_existing_target_wallet(test_client, token_value_fixture):
    CategoryFactory.create(name='transference')
    source_wallet = WalletBuilder().name('cash').create()
    response = test_client.post(
        '/transfers',
        json={
            'source': source_wallet.name,
            'target': 'target wallet',
            'amount': '100000',
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )

    assert response.status_code == 400
    assert response.headers['content-type'] == 'application/json'


def test_get(test_client, token_value_fixture):
    CategoryFactory.create(name='transference')
    wallet_source = WalletBuilder().create()
    wallet_target = WalletBuilder().name('cash').create()
    post_response = test_client.post(
        '/transfers',
        json={
            'source': wallet_source.name,
            'target': wallet_target.name,
            'amount': '100000',
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = post_response.json()
    response = test_client.get(
        f'/transfers/{response_data["uuid"]}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['source'] == wallet_source.name
    assert response_data['target'] == wallet_target.name
    assert response_data['amount'] == '100000'


def test_get_with_wrong_uuid(test_client, token_value_fixture):
    response = test_client.get(f'transfers/{uuid.uuid4()}', headers={'Authorization': f'token {token_value_fixture}'})

    assert response.status_code == 404
    assert response.headers['content-type'] == 'application/json'
