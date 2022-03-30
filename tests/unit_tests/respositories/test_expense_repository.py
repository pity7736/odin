from odin.repositories import ExpenseRepository


def test_get_expense_by_uuid(expense_fixture):
    repository = ExpenseRepository()
    got_expense = repository.get_by(uuid=expense_fixture.uuid)

    assert got_expense.date == expense_fixture.date
    assert got_expense.amount == expense_fixture.amount
    assert got_expense.uuid == expense_fixture.uuid
