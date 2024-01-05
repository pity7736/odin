from odin.accounts.application.repositories import UserRepository, TokenRepository
from odin.accounts.domain import Token
from odin.accounts.domain.crypto import get_random_string


class SessionStarter:

    def __init__(self, email: str, raw_password: str, user_repository: UserRepository,
                 token_repository: TokenRepository):
        self._email = email
        self._raw_password = raw_password
        self._user_repository = user_repository
        self._token_repository = token_repository

    async def start(self) -> Token:
        user = await self._user_repository.get_by_email(self._email)
        if user and user.check_password(self._raw_password):
            token = Token(
                value=get_random_string(length=50),
                user=user
            )
            await self._token_repository.add(token)
            return token
        else:
            raise ValueError('email or password are wrong')


class SessionFinalizer:

    def __init__(self, token_value: str, token_repository: TokenRepository):
        self._token_value = token_value
        self._token_repository = token_repository

    async def finalize(self):
        await self._token_repository.delete_by_value(self._token_value)
