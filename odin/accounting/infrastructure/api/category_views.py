from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.accounting.application.use_cases import CategoryCreator
from odin.accounting.domain import CategoryType
from odin.accounting.infrastructure.repositories import RepositoryFactory
from odin.accounts.infrastructure.api.decorators import login_required


class CategoriesEndpoint(HTTPEndpoint):

    @staticmethod
    @login_required
    async def get(request):
        try:
            category_type_name = request.query_params['type']
        except KeyError:
            return JSONResponse({'error': 'type query param is required'}, status_code=400)
        else:
            try:
                category_type = CategoryType(category_type_name)
            except ValueError:
                return JSONResponse(
                    {'error': f'type {category_type_name} is not valid category type'},
                    status_code=400
                )

        categories = []
        for category in await RepositoryFactory().get_category_repository().get_all_by_user_and_type(
                request.user,
                category_type
        ):
            categories.append({'name': category.name})
        return JSONResponse({'categories': categories})

    @staticmethod
    @login_required
    async def post(request):
        data = await request.json()
        creator = CategoryCreator(
            name=data['name'],
            type=data['type'],
            user=request.user,
            category_repository=RepositoryFactory().get_category_repository()
        )
        category = await creator.create()
        return JSONResponse({'name': category.name}, status_code=201)
