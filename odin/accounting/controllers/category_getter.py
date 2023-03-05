from odin.accounting.models import Category
from odin.accounting.repositories.repository_factory import get_category_repository


class CategoryGetter:

    def __init__(self):
        self._repository = get_category_repository()

    def get_all(self) -> tuple[Category]:
        return self._repository.get_all()

    def get_by_name(self, name: str) -> Category | None:
        return self._repository.get_by_name(name=name)
