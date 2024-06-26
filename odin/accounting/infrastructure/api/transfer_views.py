from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.application.use_cases import TransferCreator
from odin.accounting.infrastructure.repositories import RepositoryFactory
from odin.accounts.infrastructure.api.decorators import login_required


class TransfersEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    async def post(request):
        data = await request.json()
        repository_factory = RepositoryFactory()
        try:
            transfer_creator = await TransferCreator.from_wallet_ids(
                source_id=data['source'],
                target_id=data['target'],
                wallet_repository=repository_factory.get_wallet_repository(),
                transfer_repository=repository_factory.get_transfer_repository(),
                category_repository=repository_factory.get_category_repository()
            )
        except ValueError:
            return JSONResponse({}, status_code=400)
        else:
            transfer = await transfer_creator.transfer(amount=data['amount'])
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
    async def get(request):
        transfer = await RepositoryFactory().get_transfer_repository().get_by_id(request.path_params['id'])
        if transfer:
            return JSONResponse({
                'source': transfer.source.name,
                'target': transfer.target.name,
                'amount': f'{transfer.amount:f}',
                'id': transfer.id
            }, status_code=200)
        return JSONResponse({}, status_code=404)
