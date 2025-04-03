from mcp.types import Tool
from typing import Callable

ALL_TOOLS = []
TOOL_FUNC_MAP = {}

def register_tool(tool: Tool, func: Callable):
    ALL_TOOLS.append(tool)
    TOOL_FUNC_MAP[tool.name] = func

async def run(name:str, arguments:dict) -> str:
    if name not in TOOL_FUNC_MAP:
        return f"Unknown tool: {name}"
    
    return await TOOL_FUNC_MAP[name](arguments)

from . import create_black_custom_rule, create_http_application
