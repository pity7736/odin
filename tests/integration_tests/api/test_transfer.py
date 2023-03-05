from decimal import Decimal

from odin.accounting.repositories.edgedb_repositories import EdgeDBTransferRepository
from tests.factories import WalletBuilder


def test_create(test_client, token_value_fixture, transfer_category):
    wallet_source = WalletBuilder().create()
    wallet_target = WalletBuilder().name('cash').create()
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
    repository = EdgeDBTransferRepository()
    transfer = repository.get_by_uuid(response_data['uuid'])

    assert response.status_code == 201
    assert response.headers['content-type'] == 'application/json'
    assert response_data['source'] == wallet_source.name
    assert response_data['target'] == wallet_target.name
    assert response_data['amount'] == '100000'
    assert transfer.amount == Decimal('100000')


def test_get(test_client, token_value_fixture, transfer_category):
    wallet_source = WalletBuilder().create()
    wallet_target = WalletBuilder().name('cash').create()
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
        f'/accounting/transfers/{response_data["uuid"]}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['source'] == wallet_source.name
    assert response_data['target'] == wallet_target.name
    assert response_data['amount'] == '100000'
