import uuid

from nyoibo import Entity, fields

from odin.accounting.domain import Category, CategoryType
from odin.accounts.domain import User
from ..repositories import CategoryRepository


class CategoryCreator(Entity):
    _name = fields.StrField()
    _user = fields.LinkField(to=User)
    _type = fields.StrField(choices=CategoryType)

    def __init__(self, category_repository: CategoryRepository, **kwargs):
        super().__init__(**kwargs)
        self._repository = category_repository

    def create(self) -> Category:
        if self._repository.get_by_name_and_user(self._name, self._user):
            raise ValueError(f'there is already a category with name {self._name}')

        category = Category(name=self.name, id=uuid.uuid4(), user=self._user, type=self._type)
        self._repository.add(category)
        return category
