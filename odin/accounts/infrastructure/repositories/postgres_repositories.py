from typing import Optional

from odin.accounts.application.repositories import UserRepository, TokenRepository
from odin.accounts.domain import User, Token


class PostgresUserRepository(UserRepository):

    def add(self, user: User):
        pass

    def get_by_email(self, email: str) -> Optional[User]:
        pass


class PostgresTokenRepository(TokenRepository):

    def add(self, token: Token):
        pass

    def get_by_value(self, value: str) -> Optional[Token]:
        pass

    def delete_by_value(self, value: str):
        pass
