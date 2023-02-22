from starlette.responses import JSONResponse
from starlette.routing import Route

from odin.accounts.repositories import get_user_repository
from odin.auth.decorators import login_required
from odin.auth.models import Token
from odin.auth.repositories import get_token_repository
from odin.utils import get_random_string


async def login_view(request):
    data = await request.json()
    email = data['email']
    password = data['password']
    repository = get_user_repository()
    user = repository.get_by_email(email)
    if user and user.check_password(password):
        token = Token(
            value=get_random_string(length=50),
            user=user
        )
        repository = get_token_repository()
        repository.add(token)
        return JSONResponse({'token': token.value}, status_code=201)
    return JSONResponse({'message': 'email or password are wrong'}, status_code=400)


@login_required
async def logout_view(request):
    token_value = request.headers['Authorization']
    repository = get_token_repository()
    repository.delete_by_value(token_value.split()[1])
    return JSONResponse({})


routes = (
    Route('/login', login_view, methods=['POST']),
    Route('/logout', logout_view, methods=['POST'])
)
