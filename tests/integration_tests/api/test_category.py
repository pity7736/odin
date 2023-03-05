from odin.accounting.repositories.edgedb_repositories import EdgeDBCategoryRepository


def test_create_category(test_client, token_value_fixture):
    category_name = 'test category'
    response = test_client.post(
        '/accounting/categories',
        json={'name': category_name},
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()
    repository = EdgeDBCategoryRepository()
    category = repository.get_by_name(category_name)

    assert response_data['name'] == category_name
    assert category.name == category_name


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
