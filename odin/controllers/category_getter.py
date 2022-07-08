from odin.models import Category
from odin.repositories import CategoryRepository


class CategoryGetter:

    def __init__(self):
        self._repository = CategoryRepository()

    def get_all(self) -> tuple[Category]:
        return self._repository.get_all()
