from httpx import AsyncClient
from config import GLOBAL_CONFIG
import httpx

def check_slce_response(response: httpx.Response) -> str:
    if response.status_code != 200:
        return f"response status code: {response.status_code}"
    
    data = response.json()
    if data["msg"] is not None and data["msg"] != "":
        return f"request SafeLine API failed: {data['msg']}"
    
    if data["err"] is not None and data["err"] != "":
        return f"request SafeLine API failed: {data['err']}"
    
    return "success"

async def post_slce_api(path: str,req_body: dict) -> str:
    if not path.startswith("/"):
        path = f"/{path}"

    try:
        async with AsyncClient(verify=False) as client:
            response = await client.post(f"{GLOBAL_CONFIG.SAFELINE_ADDRESS}{path}", json=req_body, headers={
                "X-SLCE-API-TOKEN": f"{GLOBAL_CONFIG.SAFELINE_API_TOKEN}"
            })
            return check_slce_response(response)
    except Exception as e:
        return str(e)