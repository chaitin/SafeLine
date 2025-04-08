from httpx import AsyncClient
from config import GLOBAL_CONFIG
import httpx
import json

def get_response_data(response: httpx.Response) -> dict:
    if response.status_code != 200:
        raise Exception(f"response status code: {response.status_code}")
    
    data = response.json()
    if data["msg"] is not None and data["msg"] != "":
        raise Exception(f"request SafeLine API failed: {data['msg']}")
    
    if data["err"] is not None and data["err"] != "":
        raise Exception(f"request SafeLine API failed: {data['err']}")
    
    return data['data']

def check_slce_response(response: httpx.Response) -> str:
    try:
        get_response_data(response)
    except Exception as e:
        return str(e)
    
    return "success"


def check_slce_get_response(response: httpx.Response) -> str:
    try:
        data = get_response_data(response)
        if data:
            return json.dumps(data)
        return "empty response data"
    except Exception as e:
        return str(e)

async def get_slce_api(path: str) -> str:
    if not path.startswith("/"):
        path = f"/{path}"

    try:
        async with AsyncClient(verify=False) as client:
            response = await client.get(f"{GLOBAL_CONFIG.SAFELINE_ADDRESS}{path}", headers={
                "X-SLCE-API-TOKEN": f"{GLOBAL_CONFIG.SAFELINE_API_TOKEN}"
            })
            return check_slce_get_response(response)
    except Exception as e:
        return str(e)
    

async def post_slce_api(path: str, req_body: dict) -> str:
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