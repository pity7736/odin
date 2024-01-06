from odin.accounting.domain import CategoryType


def test_login(db_connection, user_fixture, test_client):
    login_response = test_client.post(
        '/auth/login',
        json={
            'email': user_fixture.email,
            'password': 'test'
        }
    )
    response_data = login_response.json()
    token = response_data['token']
    response = test_client.get(
        f'/accounting/categories?type={CategoryType.EXPENSE.value}',
        headers={'Authorization': f'token {token}'}
    )

    assert login_response.status_code == 201
    assert response.status_code == 200


def test_logout(db_connection, token_value_fixture, test_client):
    response = test_client.post('/auth/logout', headers={'Authorization': f'token {token_value_fixture}'})
    second_response = test_client.post('/auth/logout', headers={'Authorization': f'token {token_value_fixture}'})
    second_response_data = second_response.json()

    assert response.status_code == 200
    assert second_response.status_code == 400
    assert second_response_data['message'] == 'invalid token'
