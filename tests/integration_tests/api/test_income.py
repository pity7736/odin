import datetime

from nyoibo.fields import Decimal

from odin.accounting.repositories.edgedb_repositories import EdgeDBWalletRepository
from tests.factories import CategoryFactory, WalletBuilder


def test_create_income(test_client, token_value_fixture):
    category = CategoryFactory.create()
    wallet = WalletBuilder().create()
    response = test_client.post(
        f'/accounting/wallets/{wallet.name}/incomes',
        json={
            'date': '2022-07-31',
            'amount': '10000000',
            'category': category.name,
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()
    wallet_repository = EdgeDBWalletRepository()
    wallet = wallet_repository.get_by_name_with_incomes(wallet.name)
    income = wallet.incomes[0]

    assert response.status_code == 201
    assert response.headers['content-type'] == 'application/json'
    assert response_data['amount'] == '10000000'
    assert response_data['date'] == '2022-07-31'
    assert income.date == datetime.date(2022, 7, 31)
    assert income.amount == Decimal('10_000_000')
    assert response_data['uuid'] == income.uuid
