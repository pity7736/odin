from starlette.applications import Starlette
from starlette.routing import Route

from . import expense_views, category_views, wallet_views, income_views, transfer_views

routes = (
    Route('/categories', category_views.CategoriesEndpoint),
    Route('/transfers', transfer_views.TransfersEndpoint),
    Route('/transfers/{id}', transfer_views.TransferEndpoint),
    Route('/wallets', wallet_views.WalletsEndpoint),
    Route('/wallets/{wallet_id}', wallet_views.WalletEndpoint),
    Route('/wallets/{wallet_id}/expenses', expense_views.ExpensesEndpoint),
    Route('/wallets/{wallet_id}/expenses/{id}', expense_views.ExpenseEndpoint),
    Route('/wallets/{wallet_id}/incomes', income_views.IncomesEndpoint),
    Route('/wallets/{wallet_id}/incomes/{id}', income_views.IncomeEndpoint),
)

accounting_api = Starlette(routes=routes)
