from odin.accounting.models import Category


class InMemoryCategoryRepository:

    _categories: list[str] = []

    def add(self, category: Category):
        assert isinstance(category, Category), 'category argument must be Category instance'
        category_name = category.name.lower()
        if category_name in self._categories:
            raise ValueError(f'a category with name {category_name} already exists')
        self.__class__._categories.append(category_name)

    def get_all(self) -> tuple[Category]:
        return tuple(Category(name=name) for name in self._categories)

    def get_by_name(self, name) -> Category | None:
        if name:
            name = name.lower()
            for category_name in self._categories:
                if category_name == name:
                    return Category(name=name)
