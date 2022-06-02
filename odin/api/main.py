from starlette.applications import Starlette
from starlette.routing import Route

from . import views


routes = [
    Route('/expenses', views.Expense),
    Route('/expenses/{uuid}', views.get_expense),
    Route('/categories', views.create_category, methods=['GET', 'POST'])
]

app = Starlette(routes=routes)
