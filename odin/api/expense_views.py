from nyoibo.exceptions import RequiredValueError, FieldValueError
from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.controllers import CategoryGetter, ExpenseCreator, ExpenseGetter
from odin.repositories import WalletRepository


class ExpensesEndpoint(HTTPEndpoint):

    @staticmethod
    async def post(request):
        data = await request.json()
        category = CategoryGetter().get_by_name(data.get('category'))
        if category is None:
            return JSONResponse({}, status_code=400)

        data['category'] = category
        data['wallet'] = WalletRepository().get_by_name(data.get('wallet'))
        try:
            expense_creator = ExpenseCreator(**data)
            expense = expense_creator.create()
        except (RequiredValueError, FieldValueError, ValueError) as error:
            status_code = 400
            response_data = {'error': str(error)}
        else:
            status_code = 201
            response_data = {
                'date': expense.date.isoformat(),
                'amount': str(expense.amount),
                'uuid': expense.uuid,
                'category': category.name
            }
        return JSONResponse(response_data, status_code=status_code)

    @staticmethod
    def get(request):
        expense_getter = ExpenseGetter()
        expenses = expense_getter.all()
        serialized_expenses = []
        for expense in expenses:
            serialized_expenses.append({
                'date': expense.date.isoformat(),
                'amount': str(expense.amount),
                'uuid': expense.uuid,
                'category': expense.category.name
            })
        return JSONResponse({'expenses': serialized_expenses})


class ExpenseEndpoint(HTTPEndpoint):

    @staticmethod
    def get(request):
        expense_getter = ExpenseGetter()
        expense = expense_getter.get_by_uuid(request.path_params['uuid'])
        if expense:
            return JSONResponse(
                {
                    'date': expense.date.isoformat(),
                    'amount': str(expense.amount),
                    'uuid': expense.uuid,
                    'category': expense.category.name
                },
                status_code=200
            )
        return JSONResponse({}, status_code=404)
