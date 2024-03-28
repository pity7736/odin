from typing import Optional

import asyncpg

from odin import settings
from odin.accounts.application.repositories import UserRepository, TokenRepository
from odin.accounts.domain import User, Token
from odin.shared.db_connection_manager import DBConnectionManager


class PostgresUserRepository(UserRepository):

    def __init__(self):
        self._connection_manager = DBConnectionManager()

    async def add(self, user: User):
        async with self._connection_manager as connection:
            await connection.execute(
                'insert into users (id, email, password) values ($1, $2, $3)',
                user.id,
                user.email,
                user.password
            )

    async def get_by_email(self, email: str) -> Optional[User]:
        async with self._connection_manager as connection:
            record = await connection.fetchrow('select * from users where email = $1', email)
        if record:
            return User(**record)


class PostgresTokenRepository(TokenRepository):

    def __init__(self):
        self._connection_manager = DBConnectionManager()

    async def add(self, token: Token):
        async with self._connection_manager as connection:
            await connection.execute(
                'insert into tokens (value, user_id) VALUES ($1, $2)',
                token.value,
                token.user.id
            )

    async def get_by_value(self, value: str) -> Optional[Token]:
        async with self._connection_manager as connection:
            record = await connection.fetchrow(
                '''
                    select tokens.value, users.id, users.email, users.password
                    from tokens
                        JOIN users on (tokens.user_id = users.id)
                    where tokens.value = $1
                ''',
                value
            )

        if record:
            user = User(
                email=record['email'],
                password=record['password'],
                id=record['id']
            )
            return Token(
                value=value,
                user=user
            )

    async def delete_by_value(self, value: str):
        async with self._connection_manager as connection:
            await connection.execute('delete from tokens where value = $1', value)
