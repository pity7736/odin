
class CategoryRepository:

    _categories = {}

    def add(self, category):
        self.__class__._categories[category.name] = category

    def get_all(self):
        return tuple(self.__class__._categories.values())
