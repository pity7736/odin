from odin.controllers import ExpenseGetter


def test_get_expense_by_uuid(expense_fixture):
    expense_getter = ExpenseGetter()
    got_expense = expense_getter.get_by_uuid(uuid=expense_fixture.uuid)

    assert got_expense.date == expense_fixture.date
    assert got_expense.amount == expense_fixture.amount
