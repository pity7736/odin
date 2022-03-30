from pytest import fixture

from odin.controllers import ExpenseCreator


@fixture
def expense_fixture():
    expense_creator = ExpenseCreator(
        date='2022-03-27',
        amount='100_000'
    )
    return expense_creator.create()
