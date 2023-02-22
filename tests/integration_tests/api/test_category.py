
def test_get_all_categories(token_value_fixture, test_client):
    test_client.post(
        '/accounting/categories',
        json={'name': 'test category'},
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response = test_client.get('/accounting/categories', headers={'Authorization': f'token {token_value_fixture}'})
    response_data = response.json()

    assert response.status_code == 200
    assert response_data['categories'] == [{'name': 'test category'}]
