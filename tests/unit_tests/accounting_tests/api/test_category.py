from pytest import mark

from odin.accounting.domain import CategoryType

params = (
    ('test category0', CategoryType.EXPENSE.value),
    ('test category1', CategoryType.INCOME.value)
)


@mark.parametrize('category_name, category_type', params)
def test_create_category(category_name, category_type, test_client, token_value_fixture, category_repository):
    category_repository.get_by_name_and_user.return_value = None
    response = test_client.post(
        '/accounting/categories',
        json={'name': category_name, 'type': category_type},
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.headers['content-type'] == 'application/json'
    assert response.status_code == 201
    assert response_data['name'] == category_name


def test_get_all_category_without_type(test_client, token_value_fixture, category_repository):
    response = test_client.get(
        '/accounting/categories',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 400
    assert response_data['error'] == 'type query param is required'


def test_get_all_category_with_invalid_type(test_client, token_value_fixture, category_repository):
    response = test_client.get(
        '/accounting/categories?type=some_type',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 400
    assert response_data['error'] == 'type some_type is not valid category type'
