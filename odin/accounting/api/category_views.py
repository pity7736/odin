from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.controllers import CategoryCreator, CategoryGetter


class CategoriesEndpoint(HTTPEndpoint):

    @staticmethod
    def get(request):
        categories = []
        getter = CategoryGetter()
        for category in getter.get_all():
            categories.append({'name': category.name})
        return JSONResponse({'categories': categories})

    @staticmethod
    async def post(request):
        data = await request.json()
        creator = CategoryCreator(name=data['name'])
        category = creator.create()
        return JSONResponse({'name': category.name}, status_code=201)
