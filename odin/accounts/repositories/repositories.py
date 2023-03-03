from abc import ABCMeta, abstractmethod
from typing import Optional


from odin.accounts.models import User


class UserRepository(metaclass=ABCMeta):

    @abstractmethod
    def add(self, user: User):
        pass

    @abstractmethod
    def get_by_email(self, email: str) -> Optional[User]:
        pass
