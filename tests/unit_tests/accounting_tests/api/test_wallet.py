import uuid

from pytest import mark

from tests.factories import WalletBuilder


@mark.asyncio
async def test_create_wallet_with_existing_name(test_client, token_value_fixture, wallet_repository):
    wallet = await WalletBuilder().create()
    wallet_repository.get_by_name.return_value = wallet
    response = test_client.post(
        '/accounting/wallets',
        json={
            'name': wallet.name,
            'balance': '10000000',
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )

    assert response.status_code == 400
    assert response.headers['content-type'] == 'application/json'
    wallet_repository.get_by_name.assert_called_once_with(wallet.name)


def test_get_wallet_when_it_does_not_exists(test_client, token_value_fixture, wallet_repository):
    wallet_repository.get_by_id.return_value = None
    wallet_id = str(uuid.uuid4())
    response = test_client.get(
        f'/accounting/wallets/{wallet_id}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 404
    assert response.headers['content-type'] == 'application/json'
    assert response_data == {}
    wallet_repository.get_by_id.assert_called_once_with(wallet_id)
