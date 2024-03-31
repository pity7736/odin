from pytest import mark

from tests.factories import WalletBuilder
from tests.utils import is_uuid


@mark.asyncio
async def test_create(db_connection, test_client, token_value_fixture, transfer_category, wallet_repository,
                      transfer_repository, user_fixture):
    wallet_source = await WalletBuilder().user(user_fixture).create()
    wallet_target = await WalletBuilder().user(user_fixture).name('cash').create()
    response = test_client.post(
        '/accounting/transfers',
        json={
            'source': wallet_source.id,
            'target': wallet_target.id,
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
    assert is_uuid(response_data['id'])


@mark.asyncio
async def test_get(db_connection, test_client, token_value_fixture, transfer_category, wallet_repository,
                   transfer_repository, user_fixture):
    wallet_source = await WalletBuilder().user(user_fixture).create()
    wallet_target = await WalletBuilder().user(user_fixture).name('cash').create()
    post_response = test_client.post(
        '/accounting/transfers',
        json={
            'source': wallet_source.id,
            'target': wallet_target.id,
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
    assert response_data['amount'] == '100000.00'
