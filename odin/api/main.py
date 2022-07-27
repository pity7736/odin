from starlette.applications import Starlette
from starlette.routing import Route

from . import views


routes = [
    Route('/expenses', views.ExpensesEndpoint),
    Route('/expenses/{uuid}', views.ExpenseEndpoint),
    Route('/categories', views.CategoriesEndpoint),
    Route('/wallets', views.WalletsEndpoint),
    Route('/wallets/{name}', views.WalletEndpoint),
]

app = Starlette(routes=routes)
