import datetime
from decimal import Decimal

from odin.accounting.repositories.edgedb_repositories import EdgeDBWalletRepository
from tests.factories import CategoryFactory, WalletBuilder


def test_create(test_client, token_value_fixture):
    category = CategoryFactory.create()
    wallet = WalletBuilder().create()
    response = test_client.post(
        f'/accounting/wallets/{wallet.name}/expenses',
        json={
            'date': '2022-03-27',
            'amount': '100000',
            'category': category.name,
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()
    wallet_repository = EdgeDBWalletRepository()
    wallet = wallet_repository.get_by_name_with_expenses(wallet.name)
    expense = wallet.expenses[0]

    assert response.status_code == 201
    assert response_data['date'] == '2022-03-27'
    assert response_data['amount'] == '100000'
    assert response_data['category'] == category.name
    assert expense.date == datetime.date(2022, 3, 27)
    assert expense.amount == Decimal('100000')
    assert wallet.balance == Decimal('900_000')


def test_get_expense(test_client, token_value_fixture):
    category = CategoryFactory.create()
    wallet = WalletBuilder().create()
    post_response = test_client.post(
        f'/accounting/wallets/{wallet.name}/expenses',
        json={
            'date': '2022-03-27',
            'amount': '100000',
            'category': category.name,
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    post_response_data = post_response.json()
    response = test_client.get(
        f'/accounting/wallets/{wallet.name}/expenses/{post_response_data["uuid"]}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['date'] == post_response_data['date']
    assert response_data['amount'] == post_response_data['amount']
    assert response_data['category'] == post_response_data['category']
    assert response_data['uuid'] == post_response_data['uuid']


def test_get_all_expenses(test_client, token_value_fixture):
    wallet_builder = WalletBuilder() \
        .add_expense(amount='100000') \
        .add_expense(amount='50000') \
        .add_expense(amount='20000')
    wallet = wallet_builder.create()

    response = test_client.get(
        f'/accounting/wallets/{wallet.name}/expenses',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()
    expected_expenses = [
        {
            'date': expense.date.isoformat(),
            'amount': f'{expense.amount:f}',
            'uuid': expense.uuid,
            'category': expense.category.name
        }
        for expense in wallet.expenses
    ]
    assert response.status_code == 200
    assert response.headers['content-type'] == 'application/json'
    assert response_data['expenses'] == expected_expenses
