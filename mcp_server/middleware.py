from starlette.requests import HTTPConnection
from starlette.responses import PlainTextResponse
from starlette.types import ASGIApp, Receive, Scope, Send
from config import GLOBAL_CONFIG


class AuthenticationMiddleware:
    def __init__(
        self,
        app: ASGIApp,
    ) -> None:
        self.app = app

    async def __call__(self, scope: Scope, receive: Receive, send: Send) -> None:
        conn = HTTPConnection(scope)
        if GLOBAL_CONFIG.SECRET and GLOBAL_CONFIG.SECRET != "" and conn.headers.get("Secret") != GLOBAL_CONFIG.SECRET:
            response = PlainTextResponse("Unauthorized", status_code=401)
            await response(scope, receive, send)
            return

        await self.app(scope, receive, send)
