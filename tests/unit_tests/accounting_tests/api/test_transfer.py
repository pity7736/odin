import re
import uuid

from pytest import mark

from tests.factories import WalletBuilder
from tests.utils import UUID_PATTERN


@mark.asyncio
async def test_create(test_client, token_value_fixture, transfer_category, wallet_repository, transfer_repository,
                      user_fixture):
    wallet_source = await WalletBuilder().user(user_fixture).create()
    wallet_target = await WalletBuilder().user(user_fixture).name('cash').create()
    response = test_client.post(
        '/accounting/transfers',
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
    assert re.match(UUID_PATTERN, response_data['id'])


@mark.asyncio
async def test_create_with_non_existing_source_wallet(test_client, token_value_fixture, transfer_category):
    wallet_target = await WalletBuilder().name('cash').create()
    response = test_client.post(
        '/accounting/transfers',
        json={
            'source': 'source wallet',
            'target': wallet_target.name,
            'amount': '100000',
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )

    assert response.status_code == 400
    assert response.headers['content-type'] == 'application/json'


@mark.asyncio
async def test_create_with_non_existing_target_wallet(test_client, token_value_fixture, transfer_category):
    source_wallet = await WalletBuilder().name('cash').create()
    response = test_client.post(
        '/accounting/transfers',
        json={
            'source': source_wallet.name,
            'target': 'target wallet',
            'amount': '100000',
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )

    assert response.status_code == 400
    assert response.headers['content-type'] == 'application/json'


@mark.asyncio
async def test_get(test_client, token_value_fixture, transfer_category, wallet_repository, transfer_repository,
                   user_fixture):
    wallet_source = await WalletBuilder().user(user_fixture).create()
    wallet_target = await WalletBuilder().user(user_fixture).name('cash').create()
    post_response = test_client.post(
        '/accounting/transfers',
        json={
            'source': wallet_source.name,
            'target': wallet_target.name,
            'amount': '100000',
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = post_response.json()
    response = test_client.get(
        f'/accounting/transfers/{response_data["id"]}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['source'] == wallet_source.name
    assert response_data['target'] == wallet_target.name
    assert response_data['amount'] == '100000'


def test_get_with_wrong_uuid(test_client, token_value_fixture):
    response = test_client.get(
        f'/accounting/transfers/{uuid.uuid4()}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )

    assert response.status_code == 404
    assert response.headers['content-type'] == 'application/json'
