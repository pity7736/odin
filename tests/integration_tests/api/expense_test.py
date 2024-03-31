from pytest import mark

from tests.factories import WalletBuilder
from tests.utils import is_uuid


@mark.asyncio
async def test_create_expense(db_connection, test_client, category_fixture, wallet_fixture, token_value_fixture):
    response = test_client.post(
        f'/accounting/wallets/{wallet_fixture.id}/expenses',
        json={
            'date': '2022-03-27',
            'amount': '100000',
            'category': category_fixture.id,
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.headers['content-type'] == 'application/json'
    assert response.status_code == 201
    assert response_data['date'] == '2022-03-27'
    assert response_data['amount'] == '100000.00'
    assert response_data['category'] == category_fixture.name
    assert is_uuid(response_data['id']) is True


def test_get_expense(test_client, category_fixture, wallet_fixture, token_value_fixture):
    post_response = test_client.post(
        f'/accounting/wallets/{wallet_fixture.id}/expenses',
        json={
            'date': '2022-03-27',
            'amount': '100000',
            'category': category_fixture.id,
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    post_response_data = post_response.json()
    response = test_client.get(
        f'/accounting/wallets/{wallet_fixture.id}/expenses/{post_response_data["id"]}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['date'] == post_response_data['date']
    assert response_data['amount'] == post_response_data['amount']
    assert response_data['category'] == post_response_data['category']
    assert response_data['id'] == post_response_data['id']


@mark.asyncio
async def test_get_all_expenses(db_connection, test_client, token_value_fixture, wallet_repository, user_fixture):
    wallet_builder = WalletBuilder() \
        .user(user_fixture) \
        .add_income(amount='100000') \
        .add_expense(amount='100000') \
        .add_expense(amount='50000') \
        .add_expense(amount='20000')
    wallet = await wallet_builder.create()
    response = test_client.get(
        f'/accounting/wallets/{wallet.id}/expenses',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()
    expected_expenses = [
        {
            'date': expense.date.isoformat(),
            'amount': f'{expense.amount:.2f}',
            'id': expense.id,
            'category': expense.category.name
        }
        for expense in wallet.expenses
    ]
    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['expenses'] == expected_expenses
