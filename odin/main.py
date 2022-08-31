from starlette.applications import Starlette
from starlette.middleware import Middleware
from starlette.middleware.authentication import AuthenticationMiddleware
from starlette.routing import Route, Mount

from odin.accounting.api import category_views, expense_views, income_views, wallet_views, transference_views
from odin.auth import views
from odin.auth.backends import TokenAuthBackend, on_auth_error

routes = [
    Route('/expenses', expense_views.ExpensesEndpoint),
    Route('/expenses/{uuid}', expense_views.ExpenseEndpoint),
    Route('/categories', category_views.CategoriesEndpoint),
    Route('/wallets', wallet_views.WalletsEndpoint),
    Route('/wallets/{name}', wallet_views.WalletEndpoint),
    Route('/incomes', income_views.IncomesEndpoint),
    Route('/incomes/{uuid}', income_views.IncomeEndpoint),
    Route('/transfers', transference_views.TransfersEndpoint),
    Route('/transfers/{uuid}', transference_views.TransferenceEndpoint),
    Mount('/auth', routes=views.routes)
]

app = Starlette(
    routes=routes,
    middleware=(Middleware(AuthenticationMiddleware, backend=TokenAuthBackend(), on_error=on_auth_error),)
)
