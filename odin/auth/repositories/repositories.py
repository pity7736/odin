from abc import ABCMeta, abstractmethod
from typing import Optional

from odin.auth.models import Token


class TokenRepository(metaclass=ABCMeta):

    @abstractmethod
    def add(self, token: Token):
        pass

    @abstractmethod
    def get_by_value(self, value: str) -> Optional[Token]:
        pass
