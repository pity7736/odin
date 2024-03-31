import uuid


def test_create_income_without_category(test_client, wallet_fixture, token_value_fixture, category_repository,
                                        wallet_repository):
    category_repository.get_by_id_and_user.return_value = None
    wallet_repository.get_by_id.return_value = wallet_fixture
    response = test_client.post(
        f'/accounting/wallets/{wallet_fixture.id}/incomes',
        json={
            'date': '2022-07-31',
            'amount': '10000000',
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )

    assert response.status_code == 400


def test_get_non_existing_income(test_client, category_fixture, wallet_fixture, token_value_fixture, wallet_repository):
    wallet_repository.get_income_by_wallet_and_income_id.return_value = None
    income_id = str(uuid.uuid4())
    response = test_client.get(
        f'/accounting/wallets/{wallet_fixture.id}/incomes/{income_id}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 404
    assert response.headers['content-type'] == 'application/json'
    assert response_data == {}
    wallet_repository.get_income_by_wallet_and_income_id.assert_called_once_with(wallet_fixture.id, income_id)
