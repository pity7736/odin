from odin.models import Category


class CategoryRepository:

    _categories: dict[str, Category] = {}

    def add(self, category: Category):
        self.__class__._categories[category.name] = category

    def get_all(self) -> tuple[Category]:
        return tuple(self.__class__._categories.values())
