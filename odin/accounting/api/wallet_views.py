from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.controllers import WalletCreator
from odin.accounting.repositories.repository_factory import get_wallet_repository
from odin.auth.decorators import login_required


class WalletsEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    async def post(request):
        data = await request.json()
        repository = get_wallet_repository()
        # refactor: move this to WalletCreator
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
    @login_required
    def get(request):
        repository = get_wallet_repository()
        wallet = repository.get_by_name(request.path_params['wallet_name'])
        return JSONResponse({'name': wallet.name, 'balance': f'{wallet.balance:f}'})
