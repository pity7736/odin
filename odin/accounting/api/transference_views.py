from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.controllers import TransferenceCreator
from odin.accounting.repositories import TransferenceRepository
from odin.auth.decorators import login_required


class TransfersEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    async def post(request):
        data = await request.json()
        try:
            transference_creator = TransferenceCreator.from_wallet_names(
                source_name=data['source'],
                target_name=data['target']
            )
        except ValueError:
            return JSONResponse({}, status_code=400)
        else:
            transference = transference_creator.transfer(amount=data['amount'])
            response = {
                'source': transference.source.name,
                'target': transference.target.name,
                'amount': str(transference.amount),
                'uuid': transference.uuid
            }
            return JSONResponse(response, status_code=201)


class TransferenceEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    def get(request):
        transference = TransferenceRepository().get_by_uuid(request.path_params['uuid'])
        if transference:
            return JSONResponse({
                'source': transference.source.name,
                'target': transference.target.name,
                'amount': str(transference.amount),
                'uuid': transference.uuid
            }, status_code=200)
        return JSONResponse({}, status_code=404)
