from nyoibo import Entity, fields

from odin.accounting.models import Category
from odin.accounting.repositories import CategoryRepository


class CategoryCreator(Entity):
    _name = fields.StrField()

    def create(self) -> Category:
        category = Category(name=self.name)
        repository = CategoryRepository()
        repository.add(category)
        return category
