def test_create_wallet(db_connection, test_client, token_value_fixture):
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
    assert response_data['balance'] == '10000000.00'


def test_get_wallet(db_connection, test_client, token_value_fixture, wallet_repository):
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
        f'/accounting/wallets/{post_response_data["id"]}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['name'] == 'test wallet'
    assert response_data['balance'] == '10000000.00'
