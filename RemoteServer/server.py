from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
import requests
import os
import uvicorn

app = FastAPI()

@app.post("/mcp")
async def mcp_handler(request: Request):
    body = await request.json()

    method = body.get("method")
    req_id = body.get("id")

    # ---- initialize ----
    if method == "initialize":
        return {
            "jsonrpc": "2.0",
            "id": req_id,
            "result": {
                "capabilities": {},
                "protocolVersion": "2024-11-05"
            }
        }

    if method == "tools/list":
        return {
            "jsonrpc": "2.0",
            "id": req_id,
            "result": {
                "tools": [
                    {
                        "name": "get_pokemon",
                        "title": "Pokémon Lookup",
                        "description": "Permite consultar estadísticas, tipos y otra información detallada de un Pokémon",
                        "inputSchema": {
                            "type": "object",
                            "properties": {
                                "name": {"type": "string", "description": "Nombre del Pokémon (ej: pikachu)"}
                            },
                            "required": ["name"]
                        }
                    }
                ]
            }
        }

    if method == "tools/call":
        params = body.get("params", {})
        name = params.get("arguments", {}).get("name")
        if not name:
            return {
                "jsonrpc": "2.0",
                "id": req_id,
                "error": {"code": -32602, "message": "Missing 'name' argument"}
            }

        # Llamar a la API de PokeAPI
        resp = requests.get(f"https://pokeapi.co/api/v2/pokemon/{name.lower()}")
        if resp.status_code != 200:
            return {
                "jsonrpc": "2.0",
                "id": req_id,
                "error": {"code": 404, "message": f"Pokémon {name} no encontrado"}
            }

        data = resp.json()
        result = {
            "id": data["id"],
            "name": data["name"],
            "types": [t["type"]["name"] for t in data["types"]],
            "stats": {s["stat"]["name"]: s["base_stat"] for s in data["stats"]}
        }

        return {
            "jsonrpc": "2.0",
            "id": req_id,
            "result": {
                "content": [
                    {
                        "type": "text",
                        "text": f"Pokémon {result['name'].title()} (ID {result['id']}): "
                                f"Tipos: {', '.join(result['types'])}, "
                                f"Stats: {result['stats']}"
                    }
                ]
            }
        }

    return {
        "jsonrpc": "2.0",
        "id": req_id,
        "error": {"code": -32601, "message": f"Method {method} not found"}
    }


@app.get("/")
def root():
    return {"message": "Hello from Pokemon MCP!"}

if __name__ == "__main__":
    port = int(os.environ.get("PORT", 8080)) 
    uvicorn.run(app, host="0.0.0.0", port=port) 