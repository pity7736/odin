from starlette.applications import Starlette
from starlette.routing import Route

from . import views


routes = [
    Route('/expenses', views.Expense)
]

app = Starlette(routes=routes)
