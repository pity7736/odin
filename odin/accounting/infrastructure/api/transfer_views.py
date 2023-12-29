from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.application.use_cases import TransferCreator
from odin.accounting.infrastructure.repositories import get_transfer_repository, get_wallet_repository, \
    get_category_repository
from odin.accounts.infrastructure.api.decorators import login_required


class TransfersEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    async def post(request):
        data = await request.json()
        try:
            transfer_creator = TransferCreator.from_wallet_names(
                source_name=data['source'],
                target_name=data['target'],
                wallet_repository=get_wallet_repository(),
                transfer_repository=get_transfer_repository(),
                category_repository=get_category_repository()
            )
        except ValueError:
            return JSONResponse({}, status_code=400)
        else:
            transfer = transfer_creator.transfer(amount=data['amount'])
            response = {
                'source': transfer.source.name,
                'target': transfer.target.name,
                'amount': str(transfer.amount),
                'id': transfer.id
            }
            return JSONResponse(response, status_code=201)


class TransferEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    def get(request):
        transfer = get_transfer_repository().get_by_id(request.path_params['id'])
        if transfer:
            return JSONResponse({
                'source': transfer.source.name,
                'target': transfer.target.name,
                'amount': f'{transfer.amount:f}',
                'id': transfer.id
            }, status_code=200)
        return JSONResponse({}, status_code=404)
