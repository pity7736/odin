from nyoibo.exceptions import RequiredValueError, FieldValueError
from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.controllers import ExpenseCreator, ExpenseGetter, CategoryCreator, CategoryGetter, WalletCreator, \
    IncomeCreator
from odin.repositories import WalletRepository
from odin.repositories.income_repository import IncomeRepository


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


class CategoriesEndpoint(HTTPEndpoint):

    @staticmethod
    def get(request):
        categories = []
        getter = CategoryGetter()
        for category in getter.get_all():
            categories.append({'name': category.name})
        return JSONResponse({'categories': categories})

    @staticmethod
    async def post(request):
        data = await request.json()
        creator = CategoryCreator(name=data['name'])
        category = creator.create()
        return JSONResponse({'name': category.name}, status_code=201)


class WalletsEndpoint(HTTPEndpoint):

    @staticmethod
    async def post(request):
        data = await request.json()
        repository = WalletRepository()
        if repository.get_by_name(data['name']):
            return JSONResponse({}, status_code=400)

        wallet_creator = WalletCreator(
            name=data['name'],
            balance=data['balance']
        )
        wallet = wallet_creator.create()
        return JSONResponse({
            'name': wallet.name,
            'balance': str(wallet.balance),
            'uuid': wallet.uuid
        }, status_code=201)


class WalletEndpoint(HTTPEndpoint):

    @staticmethod
    def get(request):
        repository = WalletRepository()
        wallet = repository.get_by_name(request.path_params['name'])
        return JSONResponse({'name': wallet.name, 'balance': str(wallet.balance)})


class IncomesEndpoint(HTTPEndpoint):

    @staticmethod
    async def post(request):
        data = await request.json()
        category = CategoryGetter().get_by_name(data.get('category'))
        wallet = WalletRepository().get_by_name(data.get('wallet'))
        try:
            income_creator = IncomeCreator(
                date=data['date'],
                amount=data['amount'],
                category=category,
                wallet=wallet
            )
        except ValueError:
            return JSONResponse({}, status_code=400)

        income = income_creator.create()
        return JSONResponse({
                'date': income.date.isoformat(),
                'amount': str(income.amount),
                'uuid': str(income.uuid)
            },
            status_code=201
        )


class IncomeEndpoint(HTTPEndpoint):

    @staticmethod
    def get(request):
        income_uuid = request.path_params['uuid']
        repository = IncomeRepository()
        income = repository.get_by_uuid(uuid=income_uuid)
        return JSONResponse(
            {
                'date': income.date.isoformat(),
                'amount': str(income.amount),
                'uuid': income.uuid,
                'category': income.category.name
            }
        )
