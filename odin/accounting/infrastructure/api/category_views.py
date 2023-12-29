from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.application.use_cases import CategoryCreator
from odin.accounting.infrastructure.repositories import get_category_repository
from odin.accounts.infrastructure.api.decorators import login_required


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
        creator = CategoryCreator(name=data['name'], category_repository=get_category_repository())
        category = creator.create()
        return JSONResponse({'name': category.name}, status_code=201)
