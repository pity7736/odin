from odin import settings
from .in_memory_repositories import InMemoryTokenRepository
from .edgedb_repositories import EdgeDBTokenRepository


def get_token_repository():
    if settings.REPOSITORY == 'in-memory':
        return InMemoryTokenRepository()
    elif settings.REPOSITORY == 'edgedb':
        return EdgeDBTokenRepository()
