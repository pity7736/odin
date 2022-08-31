from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.controllers import CategoryCreator, CategoryGetter
from odin.auth.decorators import login_required


class CategoriesEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    def get(request):
        categories = []
        getter = CategoryGetter()
        for category in getter.get_all():
            categories.append({'name': category.name})
        return JSONResponse({'categories': categories})

    @staticmethod
    @login_required
    async def post(request):
        data = await request.json()
        creator = CategoryCreator(name=data['name'])
        category = creator.create()
        return JSONResponse({'name': category.name}, status_code=201)
