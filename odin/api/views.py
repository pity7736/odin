from nyoibo.exceptions import RequiredValueError, FieldValueError
from starlette.endpoints import HTTPEndpoint
from starlette.responses import JSONResponse

from odin.controllers import ExpenseCreator, ExpenseGetter, CategoryCreator, CategoryGetter


class Expense(HTTPEndpoint):

    @staticmethod
    async def post(request):
        data = await request.json()
        category = CategoryGetter().get_by_name(data.get('category'))
        if category is None:
            return JSONResponse({}, status_code=400)

        try:
            data['category'] = category
            expense_creator = ExpenseCreator(**data)
        except (RequiredValueError, FieldValueError, ValueError):
            status_code = 400
            response_data = {}
        else:
            expense = expense_creator.create()
            status_code = 201
            response_data = {
                'date': expense.date.isoformat(),
                'amount': str(expense.amount),
                'uuid': expense.uuid,
                'category': category.name
            }
        return JSONResponse(response_data, status_code=status_code)

    @staticmethod
    async def get(request):
        expense_getter = ExpenseGetter()
        expenses = expense_getter.all()
        serialized_expenses = []
        for expense in expenses:
            serialized_expenses.append({
                'date': expense.date.isoformat(),
                'amount': str(expense.amount),
                'uuid': expense.uuid,
                'category': expense.category.name
            })
        return JSONResponse({'expenses': serialized_expenses})


def get_expense(request):
    expense_getter = ExpenseGetter()
    expense = expense_getter.get_by_uuid(request.path_params['uuid'])
    if expense:
        return JSONResponse(
            {
                'date': expense.date.isoformat(),
                'amount': str(expense.amount),
                'uuid': expense.uuid,
                'category': expense.category.name
            },
            status_code=200
        )
    return JSONResponse({}, status_code=404)


async def create_category(request):
    if request.method == 'POST':
        data = await request.json()
        creator = CategoryCreator(name=data['name'])
        category = creator.create()
        return JSONResponse({'name': category.name}, status_code=201)

    categories = []
    getter = CategoryGetter()
    for category in getter.get_all():
        categories.append({'name': category.name})
    return JSONResponse({'categories': categories})
