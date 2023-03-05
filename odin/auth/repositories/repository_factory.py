from odin import settings
from .in_memory_repositories import InMemoryTokenRepository
from .edgedb_repositories import EdgeDBTokenRepository
from .repositories import TokenRepository


def get_token_repository() -> TokenRepository:
    if settings.REPOSITORY == 'in-memory':
        return InMemoryTokenRepository()
    elif settings.REPOSITORY == 'edgedb':
        return EdgeDBTokenRepository()
