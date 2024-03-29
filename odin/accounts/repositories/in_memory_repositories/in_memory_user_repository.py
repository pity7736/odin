from typing import Optional

from odin.accounts.models import User
from ..repositories import UserRepository


class InMemoryUserRepository(UserRepository):

    _user = {}

    def add(self, user):
        self.__class__._user[user.email] = User(
            email=user.email,
            password=user.password,
            first_name=user.first_name,
            last_name=user.last_name
        )

    def get_by_email(self, email) -> Optional[User]:
        return self._user.get(email)
