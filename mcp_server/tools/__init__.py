from mcp.types import Tool
from abc import ABC, abstractmethod
import os
import importlib
import logging
class ABCTool(ABC):
    @classmethod
    @abstractmethod
    async def run(self, arguments:dict) -> str:
        pass

    @classmethod
    @abstractmethod
    def tool(self) -> Tool:
        pass

class ToolRegister:
    _dict: dict[str, ABCTool] = {}
    
    @classmethod
    def register(self, tool: ABCTool) -> ABCTool:
        tool_name = tool.tool().name
        logging.info(f"Registering tool: {tool_name}")
        if tool_name in self._dict:
            raise ValueError(f"Tool {tool_name} already registered")

        self._dict[tool_name] = tool
        return tool

    def all(self) -> list[Tool]:
        return [tool.tool() for tool in self._dict.values()]
    
    async def run(self, name: str, arguments: dict) -> str:
        if name not in self._dict:
            raise ValueError(f"Unknown tool: {name}")

        return await self._dict[name].run(arguments)

def import_all_tools():
    for module in os.listdir(os.path.dirname(__file__)):
        if module == "__init__.py" or len(module) < 3 or not module.endswith(".py"):
            continue

        module_name = module[:-3]
        importlib.import_module(f".{module_name}", package=__name__)

tools = ToolRegister()

import_all_tools()