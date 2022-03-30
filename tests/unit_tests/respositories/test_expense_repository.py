from odin.repositories import ExpenseRepository


def test_get_expense_by_uuid(expense_fixture):
    repository = ExpenseRepository()
    gotten_expense = repository.get_by(uuid=expense_fixture.uuid)

    assert gotten_expense.date == expense_fixture.date
    assert gotten_expense.amount == expense_fixture.amount
    assert gotten_expense.uuid == expense_fixture.uuid
