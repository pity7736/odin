from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.controllers import TransferCreator
from odin.accounting.repositories.repository_factory import get_transfer_repository
from odin.accounts.infrastructure.api.decorators import login_required


class TransfersEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    async def post(request):
        data = await request.json()
        try:
            transfer_creator = TransferCreator.from_wallet_names(
                source_name=data['source'],
                target_name=data['target']
            )
        except ValueError:
            return JSONResponse({}, status_code=400)
        else:
            transfer = transfer_creator.transfer(amount=data['amount'])
            response = {
                'source': transfer.source.name,
                'target': transfer.target.name,
                'amount': str(transfer.amount),
                'uuid': transfer.uuid
            }
            return JSONResponse(response, status_code=201)


class TransferEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    def get(request):
        transfer = get_transfer_repository().get_by_uuid(request.path_params['uuid'])
        if transfer:
            return JSONResponse({
                'source': transfer.source.name,
                'target': transfer.target.name,
                'amount': f'{transfer.amount:f}',
                'uuid': transfer.uuid
            }, status_code=200)
        return JSONResponse({}, status_code=404)
