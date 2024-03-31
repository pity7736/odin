import uuid

from pytest import mark

from tests.factories import WalletBuilder


@mark.asyncio
async def test_create_with_non_existing_source_wallet(test_client, token_value_fixture, transfer_category,
                                                      wallet_repository):
    wallet_target = await WalletBuilder().name('cash').create()
    wallet_repository.get_by_id.side_effect = (
        None,
        wallet_target
    )
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
async def test_create_with_non_existing_target_wallet(test_client, token_value_fixture, transfer_category,
                                                      wallet_repository):
    source_wallet = await WalletBuilder().name('cash').create()
    wallet_repository.get_by_id.side_effect = (
        source_wallet,
        None
    )
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


def test_get_with_wrong_uuid(test_client, token_value_fixture, transfer_repository):
    transfer_repository.get_by_id.return_value = None
    response = test_client.get(
        f'/accounting/transfers/{uuid.uuid4()}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )

    assert response.status_code == 404
    assert response.headers['content-type'] == 'application/json'
