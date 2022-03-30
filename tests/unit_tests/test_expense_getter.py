from odin.controllers import ExpenseGetter


def test_get_expense_by_uuid(expense_fixture):
    expense_getter = ExpenseGetter()
    gotten_expense = expense_getter.get_by_uuid(uuid=expense_fixture.uuid)

    assert gotten_expense.date == expense_fixture.date
    assert gotten_expense.amount == expense_fixture.amount
