import uuid

from nyoibo import Entity, fields

from odin.accounting.domain.models import Category
from ..repositories import CategoryRepository


class CategoryCreator(Entity):
    _name = fields.StrField()

    def __init__(self, category_repository: CategoryRepository, **kwargs):
        super().__init__(**kwargs)
        self._repository = category_repository

    def create(self) -> Category:
        category = Category(name=self.name, id=uuid.uuid4())
        self._repository.add(category)
        return category
