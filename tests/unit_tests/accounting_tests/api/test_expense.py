import datetime
import uuid

from pytest import mark

from tests.factories import WalletBuilder


# TODO: test error messages text


def test_create_expense_with_wrong_category_name(test_client, category_fixture, wallet_fixture, token_value_fixture,
                                                 category_repository):
    category_repository.get_by_id_and_user.return_value = None
    category_id = 'wrong category id'
    response = test_client.post(
        f'/accounting/wallets/{wallet_fixture.id}/expenses',
        json={
            'date': '2022-03-27',
            'amount': '100000',
            'category': category_id
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )

    assert response.headers['content-type'] == 'application/json'
    assert response.status_code == 400
    category_repository.get_by_id_and_user.assert_called_once_with(category_id, category_fixture.user)


create_expense_data_params = (
    {'date': '2022-03-27'},
    {'amount': '1000'},
    {'date': 'wrong date', 'amount': '100'},
    {'date': '2022-03-29', 'amount': 'wrong amount'},
)


@mark.parametrize('data', create_expense_data_params)
def test_create_expense_with_missing_or_wrong_data(data, test_client, wallet_fixture, category_fixture,
                                                   token_value_fixture, category_repository):
    category_repository.get_by_id_and_user.return_value = category_fixture
    data['category'] = category_fixture.id
    response = test_client.post(
        f'/accounting/wallets/{wallet_fixture.id}/expenses',
        json=data,
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    assert response.status_code == 400
    category_repository.get_by_id_and_user.assert_called_once_with(category_fixture.id, category_fixture.user)


def test_create_expense_with_date_in_the_future(test_client, wallet_fixture, category_fixture, token_value_fixture,
                                                category_repository):
    category_repository.get_by_id_and_user.return_value = category_fixture
    response = test_client.post(
        f'/accounting/wallets/{wallet_fixture.name}/expenses',
        json={
            'date': (datetime.date.today() + datetime.timedelta(days=1)).isoformat(),
            'amount': '100000',
            'category': category_fixture.id
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )

    assert response.status_code == 400
    category_repository.get_by_id_and_user.assert_called_once_with(category_fixture.id, category_fixture.user)


@mark.asyncio
async def test_create_expense_with_higher_amount_that_wallet_balance(test_client, category_fixture, token_value_fixture,
                                                                     category_repository, wallet_repository):
    category_repository.get_by_id_and_user.return_value = category_fixture
    wallet = await WalletBuilder().balance('100_000').create()
    wallet_repository.get_by_id.return_value = wallet
    response = test_client.post(
        f'/accounting/wallets/{wallet.id}/expenses',
        json={
            'date': '2022-03-27',
            'amount': '102000',
            'category': category_fixture.id,
        },
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 400
    assert response_data['error'] == 'expense amount must be lower than wallet balance'
    category_repository.get_by_id_and_user.assert_called_once_with(category_fixture.id, category_fixture.user)
    wallet_repository.add.assert_called_once_with(wallet)


def test_get_non_existing_expense(test_client, wallet_fixture, token_value_fixture, wallet_repository):
    wallet_repository.get_expense_by_wallet_and_expense_id.return_value = None
    response = test_client.get(
        f'/accounting/wallets/{wallet_fixture.id}/expenses/{uuid.uuid4()}',
        headers={'Authorization': f'token {token_value_fixture}'}
    )
    response_data = response.json()

    assert response.status_code == 404
    assert response.headers['content-type'] == 'application/json'
    assert response_data == {}
