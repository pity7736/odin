from odin.accounting.domain import CategoryType


async def test_get_all_categories(db_connection, test_client, token_value_fixture):
    category_type = CategoryType.EXPENSE.value
    category_name = 'test'
    test_client.post(
        '/accounting/categories',
        json={'name': category_name, 'type': category_type},
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response = test_client.get(
        f'/accounting/categories?type={category_type}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 200
    assert response_data['categories'] == [{'name': category_name}]
