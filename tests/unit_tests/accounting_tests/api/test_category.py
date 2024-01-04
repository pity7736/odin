from pytest import mark

from tests.factories import CategoryFactory

params = (
    'test category0',
    'test category1',
)


@mark.parametrize('category_name', params)
def test_create_category(category_name, test_client, token_value_fixture):
    response = test_client.post(
        '/accounting/categories',
        json={'name': category_name},
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.headers['content-type'] == 'application/json'
    assert response.status_code == 201
    assert response_data['name'] == category_name


@mark.parametrize('category_name', params)
def test_get_all_categories(category_name, test_client, token_value_fixture, category_repository):
    CategoryFactory.create(name='category_from_some_user')
    test_client.post(
        '/accounting/categories',
        json={'name': category_name},
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response = test_client.get('/accounting/categories', headers={'Authorization': f'token {token_value_fixture}'})
    response_data = response.json()

    assert response.status_code == 200
    assert response_data['categories'] == [{'name': category_name}]
