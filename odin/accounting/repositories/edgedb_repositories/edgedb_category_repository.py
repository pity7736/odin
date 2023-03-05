from odin.accounting.models import Category

from .db_client import DBClient
from ..repositories import CategoryRepository


class EdgeDBCategoryRepository(CategoryRepository):

    _categories: list[str] = []

    def __init__(self):
        self._client = DBClient()

    def add(self, category):
        self._client.execute('insert Category { name := <str>$name }', name=category.name)

    def get_all(self) -> tuple[Category]:
        records = self._client.query('select Category {name}')
        return tuple(Category(name=record.name) for record in records)

    def get_by_name(self, name):
        if category_data := self._client.query_single('select Category {name} filter .name = <str>$name', name=name):
            return Category(name=category_data.name)
