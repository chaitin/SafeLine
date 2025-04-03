from pydantic import BaseModel, Field
from utils.request import post_slce_api
from tools import Tool, ABCTool, tools
from urllib.parse import urlparse
@tools.register
class CreateHttpApplication(BaseModel, ABCTool):
    domain: str = Field(default="",description="application domain, if empty, match all domain")
    port: int = Field(min=1, max=65535,description="application listen port, must between 1 and 65535")
    upstream: str = Field(description="application proxy address, must be a valid url")

    @classmethod
    async def run(self, arguments:dict) -> str:
        try:
            req = CreateHttpApplication.model_validate(arguments)
            parsed_upstream = urlparse(req.upstream)
            if parsed_upstream.scheme != "https" and parsed_upstream.scheme != "http":
                return "invalid upstream scheme"

            if parsed_upstream.hostname == "":
                return "invalid upstream host"
        except Exception as e:
            return str(e)

        return await post_slce_api("/api/open/site",{
                "server_names": [req.domain],
                "ports": [ str(req.port) ],
                "upstreams": [ req.upstream ],
                "type": 0,
                "static_default": 1,
                "health_check": True,
                "load_balance": {
                    "balance_type": 1
                }
            })

    @classmethod
    def tool(self) -> Tool:
        return Tool(
            name="create_http_application",
            description="在雷池 WAF 上创建一个站点应用",
            inputSchema=self.model_json_schema()
        )
