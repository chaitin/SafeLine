## mcp_server - A SafeLine WAF mcp server

- Easy to use
    - one command to run mcp_server
- Easy to develop
    - add yoor own tools to `tools` dirctory without modify other files

### quick start

```shell
docker compose -f docker-compose.yaml up -d
```

### custom your own tool

#### Hello Tool Example

This tool used to say hello to someone

1. create file `tools/hello.py`

```python
from pydantic import BaseModel, Field
from tools import Tool, ABCTool, tools

# register to global tools
@tools.register
# Hello describe function paramters
class Hello(BaseModel, ABCTool):
    # tools paramters
    name: str = Field(description="username to say hello")

    # run is tool logic, must use classmethod
    @classmethod
    async def run(arguments: dict) -> str:
        req = Hello.model_validate(arguments)
        return f"Hello {req.name}"

    # tool description, must use classmethod
    @classmethod
    def tool(self) -> Tool:
        return Tool(
            name="hello",
            description="say hello to someone",
            inputSchema=self.model_json_schema()
        )

```
