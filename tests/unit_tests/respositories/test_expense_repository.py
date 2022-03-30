from odin.controllers import ExpenseCreator
from odin.repositories import ExpenseRepository


def test_get_expense_by_uuid():
    expense_creator = ExpenseCreator(
        date='2022-03-29',
        amount='100',
    )
    expense = expense_creator.create()

    repository = ExpenseRepository()
    got_expense = repository.get_by(uuid=expense.uuid)

    assert got_expense.date == expense.date
    assert got_expense.amount == expense.amount
    assert got_expense.uuid == expense.uuid
