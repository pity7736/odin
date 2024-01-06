
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


def test_logout_without_token(user_fixture, test_client):
    response = test_client.post('/auth/logout')
    response_data = response.json()

    assert response.status_code == 401
    assert response_data['message'] == 'login required'
