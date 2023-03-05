import edgedb


class DBClient:

    _instance = None

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._client = edgedb.create_client()
        return cls._instance

    def __getattribute__(self, item):
        if item != '_client':
            return getattr(self._client, item)
        return super().__getattribute__(item)
