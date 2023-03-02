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
    assert expense.uuid == response_data['uuid']
