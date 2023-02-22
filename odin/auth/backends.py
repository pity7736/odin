import typing

from starlette.authentication import AuthenticationBackend, AuthCredentials, AuthenticationError
from starlette.requests import HTTPConnection
from starlette.responses import JSONResponse

from odin.accounts.models import User
from odin.auth.repositories import get_token_repository


def on_auth_error(request, exc):
    return JSONResponse({"message": str(exc)}, status_code=400)


class TokenAuthBackend(AuthenticationBackend):

    async def authenticate(self, conn: HTTPConnection) -> typing.Optional[tuple[AuthCredentials, User]]:
        token_value = conn.headers.get('Authorization')
        if token_value:
            repository = get_token_repository()
            token = repository.get_by_value(token_value.split()[1])
            if token:
                return AuthCredentials(), token.user
            raise AuthenticationError('invalid token')
