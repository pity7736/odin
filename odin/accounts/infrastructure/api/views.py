from starlette.responses import JSONResponse
from starlette.routing import Route

from odin.accounts.application import SessionStarter, SessionFinalizer
from odin.accounts.infrastructure.repositories import get_token_repository, get_user_repository
from odin.accounts.infrastructure.api.decorators import login_required


async def login_view(request):
    data = await request.json()
    session_starter = SessionStarter(
        email=data['email'],
        raw_password=data['password'],
        user_repository=get_user_repository(),
        token_repository=get_token_repository()
    )
    try:
        token = await session_starter.start()
    except ValueError as error:
        return JSONResponse({'message': str(error)}, status_code=400)
    else:
        return JSONResponse({'token': token.value}, status_code=201)


@login_required
async def logout_view(request):
    token_value = request.headers['Authorization']
    session_finalizer = SessionFinalizer(
        token_value=token_value.split()[1],
        token_repository=get_token_repository()
    )
    await session_finalizer.finalize()
    return JSONResponse({})


routes = (
    Route('/login', login_view, methods=['POST']),
    Route('/logout', logout_view, methods=['POST'])
)
