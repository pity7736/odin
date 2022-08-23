from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.controllers import WalletCreator
from odin.accounting.repositories import WalletRepository


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
