from starlette.endpoints import HTTPEndpoint
from starlette.requests import Request
from starlette.responses import JSONResponse

from odin.accounting.application.use_cases import IncomeCreator
from odin.accounting.infrastructure.repositories import get_wallet_repository, get_category_repository
from odin.accounts.infrastructure.api.decorators import login_required


class IncomesEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    async def post(request: Request):
        data = await request.json()
        category = await get_category_repository().get_by_name(data.get('category'))
        wallet_repository = get_wallet_repository()
        wallet = await wallet_repository.get_by_name(request.path_params['wallet_name'])
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


class IncomeEndpoint(HTTPEndpoint):

    @staticmethod
    async def get(request):
        wallet = await get_wallet_repository().get_by_name_with_incomes(request.path_params['wallet_name'])
        for income in wallet.incomes:
            if income.id == request.path_params['id']:
                return JSONResponse(
                    {
                        'date': income.date.isoformat(),
                        'amount': f'{income.amount:f}',
                        'id': income.id,
                        'category': income.category.name
                    },
                    status_code=200
                )
