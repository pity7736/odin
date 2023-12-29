from typing import Optional

from odin.accounts.application.repositories import TokenRepository, UserRepository
from odin.accounts.domain import User


class InMemoryTokenRepository(TokenRepository):

    def __init__(self):
        self._tokens = {}

    def add(self, token):
        self._tokens[token.value] = token

    def get_by_value(self, value):
        return self._tokens.get(value)

    def delete_by_value(self, value):
        self._tokens.pop(value, None)


class InMemoryUserRepository(UserRepository):

    def __init__(self):
        self._users = {}

    def add(self, user):
        self._users[user.email] = User(
            email=user.email,
            password=user.password,
            first_name=user.first_name,
            last_name=user.last_name
        )

    def get_by_email(self, email) -> Optional[User]:
        return self._users.get(email)
