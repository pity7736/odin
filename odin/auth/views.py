from starlette.responses import JSONResponse
from starlette.routing import Route

from odin.accounts.repositories import UserRepository
from odin.auth.decorators import login_required
from odin.auth.models import Token
from odin.auth.repositories import TokenRepository
from odin.utils import get_random_string


async def login_view(request):
    data = await request.json()
    email = data['email']
    password = data['password']
    user = UserRepository().get_by_email(email)
    if user and user.check_password(password):
        token = Token(
            value=get_random_string(length=50),
            user=user
        )
        repository = TokenRepository()
        repository.add(token)
        return JSONResponse({'token': token.value}, status_code=201)
    return JSONResponse({'message': 'email or password are wrong'}, status_code=400)


@login_required
async def logout_view(request):
    token_value = request.headers['Authorization']
    repository = TokenRepository()
    repository.delete_by_value(token_value.split()[1])
    return JSONResponse({})


routes = (
    Route('/login', login_view, methods=['POST']),
    Route('/logout', logout_view, methods=['POST'])
)
