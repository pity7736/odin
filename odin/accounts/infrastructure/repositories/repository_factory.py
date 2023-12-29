from odin.accounts.application.repositories import UserRepository, TokenRepository
from .postgres_repositories import PostgresUserRepository, PostgresTokenRepository


def get_user_repository() -> UserRepository:
    return PostgresUserRepository()


def get_token_repository() -> TokenRepository:
    return PostgresTokenRepository()
