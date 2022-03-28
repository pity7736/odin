from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.controllers import ExpenseCreator


class Expense(HTTPEndpoint):

    @staticmethod
    async def post(request):
        data = await request.json()
        expense_creator = ExpenseCreator(**data)
        expense = expense_creator.create()
        return JSONResponse({
            'date': expense.date.isoformat(),
            'amount': str(expense.amount),
            'uuid': expense.uuid
        }, status_code=201)
