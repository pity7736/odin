from odin.models import Category


class CategoryRepository:

    _categories: dict[str, Category] = {}

    def add(self, category: Category):
        assert isinstance(category, Category), 'category argument must be Category instance'
        category_name = category.name.lower()
        if category_name in self._categories.keys():
            raise ValueError(f'a category with name {category_name} already exists')
        self.__class__._categories[category_name] = Category(name=category_name)

    def get_all(self) -> tuple[Category]:
        return tuple(self._categories.values())

    def get_by_name(self, name) -> Category | None:
        if name:
            return self._categories.get(name.lower())
