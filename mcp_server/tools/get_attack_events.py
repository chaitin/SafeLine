from pydantic import BaseModel, Field
from utils.request import get_slce_api
from tools import Tool, ABCTool, tools
from urllib.parse import urlparse
@tools.register
class CreateHttpApplication(BaseModel, ABCTool):
    ip: str = Field(default="", description="the attacker's client IP address")
    size: int = Field(default=10, min=1, max=100, description="the number of results to return")
    start: str = Field(default="", description="start time, millisecond timestamp")
    end: str = Field(default="", description="end time, millisecond timestamp")

    @classmethod
    async def run(self, arguments:dict) -> str:
        try:
            req = CreateHttpApplication.model_validate(arguments)
            parsed_upstream = urlparse(req.upstream)
            if parsed_upstream.scheme not in ["https", "http"]:
                return "invalid upstream scheme"

            if parsed_upstream.hostname == "":
                return "invalid upstream host"
        except Exception as e:
            return str(e)

        return await get_slce_api(f"api/open/events?page=1&page_size={req.size}&ip={req.ip}&start={req.start}&end={req.end}")

    @classmethod
    def tool(self) -> Tool:
        return Tool(
            name="waf_get_attack_events",
            description="获取雷池 WAF 所记录的攻击事件",
            inputSchema=self.model_json_schema()
        )
