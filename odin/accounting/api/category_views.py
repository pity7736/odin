from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.controllers import CategoryCreator
from odin.accounting.repositories.repository_factory import get_category_repository
from odin.auth.decorators import login_required


class CategoriesEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    def get(request):
        categories = []
        for category in get_category_repository().get_all():
            categories.append({'name': category.name})
        return JSONResponse({'categories': categories})

    @staticmethod
    @login_required
    async def post(request):
        data = await request.json()
        creator = CategoryCreator(name=data['name'])
        category = creator.create()
        return JSONResponse({'name': category.name}, status_code=201)
