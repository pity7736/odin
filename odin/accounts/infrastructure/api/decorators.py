import asyncio

from starlette.responses import JSONResponse

from odin.accounts.domain.models import User


def login_required(function):
    async def wrapper(request):
        if isinstance(request.user, User):
            if asyncio.iscoroutinefunction(function):
                return await function(request)
            return function(request)

        return JSONResponse({'message': 'login required'}, status_code=401)
    return wrapper
