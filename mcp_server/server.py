from mcp.server import Server
from mcp.types import Tool, TextContent
import mcp.server.sse
from starlette.applications import Starlette
from starlette.routing import Route, Mount
from starlette.requests import Request
import uvicorn
from tools import tools
from config import GLOBAL_CONFIG
from middleware import AuthenticationMiddleware
from starlette.middleware import Middleware

# Create an MCP server
mcp_server = Server("SafeLine WAF mcp server")
sse = mcp.server.sse.SseServerTransport("/messages/")

@mcp_server.list_tools()
async def list_tools() -> list[Tool]:
    return tools.all()

@mcp_server.call_tool()
async def call_tool(name: str, arguments: dict) -> list[TextContent]:
    result = await tools.run(name, arguments)
    
    return [TextContent(
        type="text",
        text=f"{name} result: {result}"
    )]



async def handle_sse(request: Request) -> None:
    async with sse.connect_sse(
        request.scope, request.receive, request._send
    ) as [read_stream, write_stream]:
        await mcp_server.run(
            read_stream, write_stream, mcp_server.create_initialization_options()
        )

def main():
    starlette_app = Starlette(debug=GLOBAL_CONFIG.DEBUG,routes=[
        Route("/sse", endpoint=handle_sse),
        Mount("/messages/", app=sse.handle_post_message, middleware=[Middleware(AuthenticationMiddleware)]),
    ])

    uvicorn.run(starlette_app, host=GLOBAL_CONFIG.LISTEN_ADDRESS, port=GLOBAL_CONFIG.LISTEN_PORT)
