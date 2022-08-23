from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.controllers import CategoryGetter, IncomeCreator
from odin.accounting.repositories import WalletRepository, IncomeRepository


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
