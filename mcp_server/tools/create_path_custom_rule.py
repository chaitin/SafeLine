from pydantic import BaseModel, Field
from utils.request import post_slce_api
from tools import Tool, ABCTool, tools

@tools.register
class CreatePathCustomRule(BaseModel, ABCTool):
    path: str = Field(description="request path to block or allow")
    action: int = Field(min=0, max=1,description="1: block, 0: allow")

    @classmethod
    async def run(self, arguments:dict) -> str:
        try:
            req = CreatePathCustomRule.model_validate(arguments)
        except Exception as e:
            return str(e)

        name = ""
        match req.action:
            case 0:
                name += "allow "
            case 1:
                name += "block "
            case _:
                return "invalid action"

        if not req.path or req.path == "":
            return "path is required"

        name += f"path: {req.path}"

        return await post_slce_api("/api/open/policy",{
            "name": name,
            "is_enabled": True,
            "pattern": [
                [
                    {
                        "k": "uri_no_query",
                        "op": "eq",
                        "v": [req.path],
                        "sub_k": ""
                    },
                ]
            ],
            "action": req.action
        })
    
    @classmethod
    def tool(self) -> Tool:
        return Tool(
            name="create_path_custom_rule",
            description="在雷池 WAF 上创建一个 url 路径的自定义黑名单或者自定义白名单",
            inputSchema=self.model_json_schema()
        )