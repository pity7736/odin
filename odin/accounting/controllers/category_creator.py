from nyoibo import Entity, fields

from odin.accounting.models import Category
from odin.accounting.repositories.repository_factory import get_category_repository


class CategoryCreator(Entity):
    _name = fields.StrField()

    def create(self) -> Category:
        category = Category(name=self.name)
        repository = get_category_repository()
        repository.add(category)
        return category
