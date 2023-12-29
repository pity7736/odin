from pytest import fixture

from tests.factories import ExpenseFactory


@fixture
def expense_fixture():
    return ExpenseFactory.create(
        date='2022-03-27',
        amount='100_00'
    )
