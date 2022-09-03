from odin.accounting.models import Category

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


class EdgeDBCategoryRepository:

    _categories: list[str] = []

    def __init__(self):
        self._client = DBClient()

    def add(self, category: Category):
        self._client.execute('insert Category { name := <str>$name }', name=category.name)

    def get_all(self) -> tuple[Category]:
        return tuple(Category(name=name) for name in self._categories)

    def get_by_name(self, name) -> Category | None:
        if category_data := self._client.query_single('select Category {name} filter .name = <str>$name', name=name):
            return Category(name=category_data.name)
