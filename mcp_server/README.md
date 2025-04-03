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
from tools import register_tool, Tool

# Hello describe function paramters
class Hello(BaseModel):
    name: str = Field(description="username to say hello")

# hello is tool logic
async def hello(arguments: dict) -> str:
    """
    Say hello to someone
    """
    return f"Hello {arguments['name']}"

# register tool to global variable
register_tool(
    Tool(
        name="hello",
        description="say hello to someone",
        inputSchema=Hello.model_json_schema()
    ),
    hello
)
```

2. import this tool in `tools/__init__.py`

```python
from . import hello
```