from starlette.applications import Starlette
from starlette.routing import Route

from . import category_views, expense_views, income_views, wallet_views


routes = [
    Route('/expenses', expense_views.ExpensesEndpoint),
    Route('/expenses/{uuid}', expense_views.ExpenseEndpoint),
    Route('/categories', category_views.CategoriesEndpoint),
    Route('/wallets', wallet_views.WalletsEndpoint),
    Route('/wallets/{name}', wallet_views.WalletEndpoint),
    Route('/incomes', income_views.IncomesEndpoint),
    Route('/incomes/{uuid}', income_views.IncomeEndpoint)
]

app = Starlette(routes=routes)
