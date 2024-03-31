from starlette.endpoints import HTTPEndpoint
from starlette.requests import Request
from starlette.responses import JSONResponse

from odin.accounting.application.use_cases import IncomeCreator
from odin.accounting.infrastructure.repositories import RepositoryFactory
from odin.accounts.infrastructure.api.decorators import login_required


class IncomesEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    async def post(request: Request):
        data = await request.json()
        repository_factory = RepositoryFactory()
        category = await repository_factory.get_category_repository().get_by_id_and_user(
            data.get('category'),
            request.user
        )
        wallet_repository = repository_factory.get_wallet_repository()
        wallet = await wallet_repository.get_by_id(request.path_params['wallet_id'])
        try:
            income_creator = IncomeCreator(
                date=data['date'],
                amount=data['amount'],
                category=category,
                wallet=wallet,
                wallet_repository=wallet_repository
            )
        except ValueError:
            return JSONResponse({}, status_code=400)

        income = await income_creator.create()
        return JSONResponse({
                'date': income.date.isoformat(),
                'amount': str(income.amount),
                'id': str(income.id)
            },
            status_code=201
        )

    @staticmethod
    @login_required
    async def get(request):
        incomes = await RepositoryFactory().get_wallet_repository().get_incomes_by_wallet_id(
            request.path_params['wallet_id'],
        )
        serialized_incomes = []
        for income in incomes:
            serialized_incomes.append({
                'date': income.date.isoformat(),
                'amount': f'{income.amount:.2f}',
                'id': income.id,
                'category': income.category.name
            })
        return JSONResponse({'incomes': serialized_incomes})


class IncomeEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    async def get(request):
        income = await RepositoryFactory().get_wallet_repository().get_income_by_wallet_and_income_id(
            request.path_params['wallet_id'],
            request.path_params['id']
        )
        if income:
            return JSONResponse(
                {
                    'date': income.date.isoformat(),
                    'amount': f'{income.amount:f}',
                    'id': income.id,
                    'category': income.category.name
                },
                status_code=200
            )
        return JSONResponse({}, status_code=404)
