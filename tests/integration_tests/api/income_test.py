from pytest import mark

from tests.factories import WalletBuilder
from tests.utils import is_uuid


def test_create_income(db_connection, test_client, category_fixture, wallet_fixture, token_value_fixture):
    response = test_client.post(
        f'/accounting/wallets/{wallet_fixture.id}/incomes',
        json={
            'date': '2022-07-31',
            'amount': '10000000',
            'category': category_fixture.id,
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 201
    assert response.headers['content-type'] == 'application/json'
    assert response_data['amount'] == '10000000'
    assert response_data['date'] == '2022-07-31'
    assert is_uuid(response_data['id'])


@mark.asyncio
async def test_get_all_expenses(db_connection, test_client, token_value_fixture, wallet_repository, user_fixture):
    wallet_builder = WalletBuilder() \
        .user(user_fixture) \
        .add_expense(amount='100000') \
        .add_income(amount='100000') \
        .add_income(amount='50000') \
        .add_income(amount='20000')
    wallet = await wallet_builder.create()
    response = test_client.get(
        f'/accounting/wallets/{wallet.id}/incomes',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()
    expected_incomes = [
        {
            'date': income.date.isoformat(),
            'amount': f'{income.amount:.2f}',
            'id': income.id,
            'category': income.category.name
        }
        for income in wallet.incomes
    ]
    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['incomes'] == expected_incomes
