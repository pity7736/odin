from nyoibo.exceptions import RequiredValueError, FieldValueError
from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.controllers import CategoryGetter, ExpenseCreator
from odin.accounting.repositories.repository_factory import get_wallet_repository
from odin.auth.decorators import login_required


class ExpensesEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    async def post(request):
        data = await request.json()
        category = CategoryGetter().get_by_name(data.get('category'))
        if category is None:
            return JSONResponse({}, status_code=400)

        data['category'] = category
        data['wallet'] = get_wallet_repository().get_by_name(request.path_params['wallet_name'])
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
    @login_required
    def get(request):
        wallet = get_wallet_repository().get_by_name_with_expenses(request.path_params['wallet_name'])
        serialized_expenses = []
        for expense in wallet.expenses:
            serialized_expenses.append({
                'date': expense.date.isoformat(),
                'amount': f'{expense.amount:f}',
                'uuid': expense.uuid,
                'category': expense.category.name
            })
        return JSONResponse({'expenses': serialized_expenses})


class ExpenseEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    def get(request):
        wallet = get_wallet_repository().get_by_name_with_expenses(request.path_params['wallet_name'])
        for expense in wallet.expenses:
            if expense.uuid == request.path_params['uuid']:
                return JSONResponse(
                    {
                        'date': expense.date.isoformat(),
                        'amount': f'{expense.amount:f}',
                        'uuid': expense.uuid,
                        'category': expense.category.name
                    },
                    status_code=200
                )
        return JSONResponse({}, status_code=404)
