from starlette.applications import Starlette
from starlette.routing import Route

from . import expense_views, category_views, wallet_views, income_views, transfer_views

routes = (
    Route('/categories', category_views.CategoriesEndpoint),
    Route('/transfers', transfer_views.TransfersEndpoint),
    Route('/transfers/{uuid}', transfer_views.TransferEndpoint),
    Route('/wallets', wallet_views.WalletsEndpoint),
    Route('/wallets/{wallet_name}', wallet_views.WalletEndpoint),
    Route('/wallets/{wallet_name}/expenses', expense_views.ExpensesEndpoint),
    Route('/wallets/{wallet_name}/expenses/{uuid}', expense_views.ExpenseEndpoint),
    Route('/wallets/{wallet_name}/incomes', income_views.IncomesEndpoint),
    Route('/wallets/{wallet_name}/incomes/{uuid}', income_views.IncomeEndpoint),
)

accounting_api = Starlette(routes=routes)
