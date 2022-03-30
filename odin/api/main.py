from starlette.applications import Starlette
from starlette.routing import Route

from . import views


routes = [
    Route('/expenses', views.Expense),
    Route('/expenses/{uuid}', views.get_expense)
]

app = Starlette(routes=routes)
