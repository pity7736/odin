from odin import settings
from .in_memory_repositories import InMemoryUserRepository
from .edgedb_repositories import EdgeDBUserRepository


def get_user_repository():
    if settings.REPOSITORY == 'in-memory':
        return InMemoryUserRepository()
    elif settings.REPOSITORY == 'edgedb':
        return EdgeDBUserRepository()
