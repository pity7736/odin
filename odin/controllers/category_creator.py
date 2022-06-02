from nyoibo import Entity, fields

from odin.models import Category
from odin.repositories import CategoryRepository


class CategoryCreator(Entity):
    _name = fields.StrField()

    def create(self):
        category = Category(name=self.name)
        repository = CategoryRepository()
        repository.add(category)
        return category
