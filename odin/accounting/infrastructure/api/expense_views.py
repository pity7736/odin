from nyoibo.exceptions import RequiredValueError, FieldValueError
from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.application.use_cases import ExpenseCreator
from odin.accounting.infrastructure.repositories import RepositoryFactory
from odin.accounts.infrastructure.api.decorators import login_required


class ExpensesEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    async def post(request):
        data = await request.json()
        repository_factory = RepositoryFactory()
        category = await repository_factory.get_category_repository().get_by_id_and_user(
            data.get('category'),
            request.user
        )
        if category is None:
            return JSONResponse({}, status_code=400)

        data['category'] = category
        wallet_repository = repository_factory.get_wallet_repository()
        data['wallet'] = await wallet_repository.get_by_id(request.path_params['wallet_id'])
        try:
            expense_creator = ExpenseCreator(**data, wallet_repository=wallet_repository)
            expense = await expense_creator.create()
        except (RequiredValueError, FieldValueError, ValueError) as error:
            status_code = 400
            response_data = {'error': str(error)}
        else:
            status_code = 201
            response_data = {
                'date': expense.date.isoformat(),
                'amount': f'{expense.amount:.2f}',
                'id': expense.id,
                'category': category.name
            }
        return JSONResponse(response_data, status_code=status_code)

    @staticmethod
    @login_required
    async def get(request):
        expenses = await RepositoryFactory().get_wallet_repository().get_expenses_by_wallet_id(
            request.path_params['wallet_id']
        )
        serialized_expenses = []
        for expense in expenses:
            serialized_expenses.append({
                'date': expense.date.isoformat(),
                'amount': f'{expense.amount:.2f}',
                'id': expense.id,
                'category': expense.category.name
            })
        return JSONResponse({'expenses': serialized_expenses})


class ExpenseEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    async def get(request):
        expense = await RepositoryFactory().get_wallet_repository().get_expense_by_wallet_and_expense_id(
            request.path_params['wallet_id'],
            request.path_params['id']
        )
        if expense:
            return JSONResponse(
                {
                    'date': expense.date.isoformat(),
                    'amount': f'{expense.amount:f}',
                    'id': expense.id,
                    'category': expense.category.name
                },
                status_code=200
            )
        return JSONResponse({}, status_code=404)
