from typing import Optional

import asyncpg

from odin import settings
from odin.accounts.application.repositories import UserRepository, TokenRepository
from odin.accounts.domain import User, Token


class PostgresUserRepository(UserRepository):

    async def add(self, user: User):
        connection = await asyncpg.connect(
            host=settings.DB_HOST,
            user=settings.DB_USER,
            database=settings.DB_NAME,
            password=settings.DB_PASSWORD,
            port=settings.DB_PORT,
        )
        await connection.execute(
            'insert into users (id, email, password) values ($1, $2, $3)',
            user.id,
            user.email,
            user.password
        )
        await connection.close()

    async def get_by_email(self, email: str) -> Optional[User]:
        connection = await asyncpg.connect(
            host=settings.DB_HOST,
            user=settings.DB_USER,
            database=settings.DB_NAME,
            password=settings.DB_PASSWORD,
            port=settings.DB_PORT,
        )
        record = await connection.fetchrow('select * from users where email = $1', email)
        await connection.close()
        if record:
            return User(**record)


class PostgresTokenRepository(TokenRepository):

    async def add(self, token: Token):
        connection = await asyncpg.connect(
            host=settings.DB_HOST,
            user=settings.DB_USER,
            database=settings.DB_NAME,
            password=settings.DB_PASSWORD,
            port=settings.DB_PORT,
        )
        await connection.execute(
            'insert into tokens (value, user_id) VALUES ($1, $2)',
            token.value,
            token.user.id
        )
        await connection.close()

    async def get_by_value(self, value: str) -> Optional[Token]:
        connection = await asyncpg.connect(
            host=settings.DB_HOST,
            user=settings.DB_USER,
            database=settings.DB_NAME,
            password=settings.DB_PASSWORD,
            port=settings.DB_PORT,
        )
        record = await connection.fetchrow(
            '''
                select tokens.value, users.id, users.email, users.password
                from tokens
                    JOIN users on (tokens.user_id = users.id)
                where tokens.value = $1
            ''',
            value
        )
        await connection.close()
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
        connection = await asyncpg.connect(
            host=settings.DB_HOST,
            user=settings.DB_USER,
            database=settings.DB_NAME,
            password=settings.DB_PASSWORD,
            port=settings.DB_PORT,
        )
        await connection.execute('delete from tokens where value = $1', value)
        await connection.close()
