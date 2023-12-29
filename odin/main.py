from starlette.applications import Starlette
from starlette.middleware import Middleware
from starlette.middleware.authentication import AuthenticationMiddleware
from starlette.routing import Mount

from odin.accounting.api import accounting_api
from odin.accounts.infrastructure.api import views
from odin.accounts.infrastructure.api.backends import TokenAuthBackend, on_auth_error

routes = [
    Mount('/accounting', routes=accounting_api.routes),
    Mount('/auth', routes=views.routes)
]

app = Starlette(
    routes=routes,
    middleware=(Middleware(AuthenticationMiddleware, backend=TokenAuthBackend(), on_error=on_auth_error),)
)
