from pydantic import BaseModel, Field
from utils.request import post_slce_api
from tools import Tool, ABCTool, tools
import ipaddress

@tools.register
class CreateIPCustomRule(BaseModel, ABCTool):
    ip: str = Field(description="request ip to allow or block")
    action: int = Field(min=0, max=1,description="1: block, 0: allow")

    @classmethod
    async def run(self, arguments:dict) -> str:
        try:
            req = CreateIPCustomRule.model_validate(arguments)
            ipaddress.ip_address(req.ip)
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

        if not req.ip or req.ip == "":
            return "ip is required"

        name += f"ip: {req.ip}"

        return await post_slce_api("/api/open/policy",{
            "name": name,
            "is_enabled": True,
            "pattern": [
                [
                    {
                        "k": "src_ip",
                        "op": "eq",
                        "v": [req.ip],
                        "sub_k": ""
                    },
                ]
            ],
            "action": req.action
        })
    
    @classmethod
    def tool(self) -> Tool:
        return Tool(
            name="create_ip_custom_rule",
            description="在雷池 WAF 上创建一个 ip 的自定义黑名单或者自定义白名单",
            inputSchema=self.model_json_schema()
        )