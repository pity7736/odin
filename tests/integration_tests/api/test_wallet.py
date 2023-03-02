from odin.accounting.repositories.edgedb_repositories import EdgeDBWalletRepository


def test_create_wallet(test_client, token_value_fixture):
    wallet_name = 'test wallet'
    response = test_client.post(
        '/accounting/wallets',
        json={
            'name': wallet_name,
            'balance': '10000000'
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()
    repository = EdgeDBWalletRepository()
    wallet = repository.get_by_name(wallet_name)

    assert response_data['name'] == wallet_name
    assert wallet.name == wallet_name


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
