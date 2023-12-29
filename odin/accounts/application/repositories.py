from abc import ABCMeta, abstractmethod
from typing import Optional


from odin.accounts.domain import User, Token


class UserRepository(metaclass=ABCMeta):

    @abstractmethod
    def add(self, user: User):
        pass

    @abstractmethod
    def get_by_email(self, email: str) -> Optional[User]:
        pass


class TokenRepository(metaclass=ABCMeta):

    @abstractmethod
    def add(self, token: Token):
        pass

    @abstractmethod
    def get_by_value(self, value: str) -> Optional[Token]:
        pass

    @abstractmethod
    def delete_by_value(self, value: str):
        pass
