from ..repositories import TokenRepository


class InMemoryTokenRepository(TokenRepository):
    _tokens = {}

    def add(self, token):
        self.__class__._tokens[token.value] = token

    def get_by_value(self, value):
        return self._tokens.get(value)

    def delete_by_value(self, value):
        self.__class__._tokens.pop(value, None)
