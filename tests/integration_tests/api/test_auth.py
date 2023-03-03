
def test_logout(token_value_fixture, test_client):
    response = test_client.post('/auth/logout', headers={'Authorization': f'token {token_value_fixture}'})
    second_response = test_client.post('/auth/logout', headers={'Authorization': f'token {token_value_fixture}'})
    second_response_data = second_response.json()

    assert response.status_code == 200
    assert second_response.status_code == 400
    assert second_response_data['message'] == 'invalid token'
