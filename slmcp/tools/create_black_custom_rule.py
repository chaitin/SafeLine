from pydantic import BaseModel, Field
from utils.request import post_slce_api
from tools import register_tool, Tool

class CreateBlackCustomRule(BaseModel):
    path: str = Field(default="",description="request path to block, if path is empty, block all paths")
    ip: str = Field(default="",description="request ip to block, if ip is empty, block all ips")

async def create_black_custom_rule(arguments:dict) -> str:
    """
    Create a new black custom rule.

    Args:
        path: request path to block, if path is empty, block all paths
        ip: request ip to block, if ip is empty, block all ips
    """

    pattern = []
    name = "block "

    path = arguments["path"]
    ip = arguments["ip"]

    if path is not None and path != "":
        pattern.append({
            "k": "uri_no_query",
            "op": "eq",
            "v": [path],
            "sub_k": ""
        })
        name += f"path: {path} "

    if ip is not None and ip != "":
        pattern.append({
            "k": "src_ip",
            "op": "eq",
            "v": [ip],
            "sub_k": ""
        })
        name += f"ip: {ip}"
    
    if len(pattern) == 0:
        return "path and ip cannot be empty"

    return await post_slce_api("/api/open/policy",{
        "name": name,
        "is_enabled": True,
        "pattern": [pattern],
        "action": 1
    })
    

register_tool(
    Tool(
        name="create_black_custom_rule",
        description="在雷池 WAF 上创建一个黑名单自定义规则",
        inputSchema=CreateBlackCustomRule.model_json_schema()
    ),
    create_black_custom_rule
)