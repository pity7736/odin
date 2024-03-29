import asyncpg

from odin import settings


class DBConnectionManager:

    def __init__(self):
        self._connection: asyncpg.Connection = None

    async def __aenter__(self) -> asyncpg.Connection:
        self._connection = await asyncpg.connect(
            host=settings.DB_HOST,
            user=settings.DB_USER,
            database=settings.DB_NAME,
            password=settings.DB_PASSWORD,
            port=settings.DB_PORT,
        )
        return self._connection

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        await self._connection.close()
