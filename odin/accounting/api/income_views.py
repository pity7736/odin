from starlette.endpoints import HTTPEndpoint
from starlette.requests import Request
from starlette.responses import JSONResponse

from odin.accounting.controllers import IncomeCreator
from odin.accounting.repositories.repository_factory import get_wallet_repository, get_category_repository
from odin.auth.decorators import login_required


class IncomesEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    async def post(request: Request):
        data = await request.json()
        category = get_category_repository().get_by_name(data.get('category'))
        wallet = get_wallet_repository().get_by_name(request.path_params['wallet_name'])
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
        wallet = get_wallet_repository().get_by_name_with_incomes(request.path_params['wallet_name'])
        for income in wallet.expenses:
            if income.uuid == request.path_params['uuid']:
                return JSONResponse(
                    {
                        'date': income.date.isoformat(),
                        'amount': f'{income.amount:f}',
                        'uuid': income.uuid,
                        'category': income.category.name
                    },
                    status_code=200
                )
