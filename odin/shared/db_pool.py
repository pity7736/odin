import asyncpg

from odin import settings


__pool = None


async def initialize_pool(min_size=1, max_size=10):
    global __pool
    if __pool is None:
        __pool = await asyncpg.create_pool(
            host=settings.DB_HOST,
            user=settings.DB_USER,
            database=settings.DB_NAME,
            password=settings.DB_PASSWORD,
            port=settings.DB_PORT,
            min_size=min_size,
            max_size=max_size,
            timeout=5
        )


def get_pool() -> asyncpg.Pool:
    if __pool is None:
        raise RuntimeError('uninitialized pool. call initialize_pool function first')
    return __pool
