
def test_login(user_fixture, test_client):
    response = test_client.post(
        '/auth/login',
        json={
            'email': user_fixture.email,
            'password': 'test'
        }
    )
    response_data = response.json()

    assert response.status_code == 201
    assert response_data['token']


def test_login_with_wrong_password(user_fixture, test_client):
    response = test_client.post(
        '/auth/login',
        json={
            'email': user_fixture.email,
            'password': 'test2'
        }
    )
    response_data = response.json()

    assert response.status_code == 400
    assert response_data['message'] == 'email or password are wrong'


def test_logout(user_fixture, test_client):
    login_response = test_client.post(
        '/auth/login',
        json={
            'email': user_fixture.email,
            'password': 'test'
        }
    )
    token = login_response.json()['token']
    response = test_client.post('/auth/logout', headers={'Authorization': f'token {token}'})
    second_response = test_client.post('/auth/logout', headers={'Authorization': f'token {token}'})
    second_response_data = second_response.json()

    assert response.status_code == 200
    assert second_response.status_code == 400
    assert second_response_data['message'] == 'invalid token'


def test_logout_with_token(user_fixture, test_client):
    response = test_client.post('/auth/logout')
    response_data = response.json()

    assert response.status_code == 401
    assert response_data['message'] == 'login required'
