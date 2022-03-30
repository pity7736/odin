from nyoibo.exceptions import RequiredValueError, FieldValueError
from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.controllers import ExpenseCreator, ExpenseGetter


class Expense(HTTPEndpoint):

    @staticmethod
    async def post(request):
        data = await request.json()
        try:
            expense_creator = ExpenseCreator(**data)
        except (RequiredValueError, FieldValueError):
            status_code = 400
            response_data = {}
        else:
            expense = expense_creator.create()
            status_code = 201
            response_data = {
                'date': expense.date.isoformat(),
                'amount': str(expense.amount),
                'uuid': expense.uuid
            }
        return JSONResponse(response_data, status_code=status_code)


def get_expense(request):
    expense_getter = ExpenseGetter()
    expense = expense_getter.get_by_uuid(request.path_params['uuid'])
    if expense:
        return JSONResponse(
            {
                'date': expense.date.isoformat(),
                'amount': str(expense.amount),
                'uuid': expense.uuid
            },
            status_code=200
        )
    return JSONResponse({}, status_code=404)
