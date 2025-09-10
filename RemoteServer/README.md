# 🐁 Pokémon (hola demonio) MCP

This project implements an MCP server that exposes the **PokeAPI** as a remote tool for integration with a chatbot.

---

## ☁️ Deployment

The MCP server is available at the endpoint `/mcp`:  
**https://redes-proyecto-mcp-673347929287.us-central1.run.app**

Use it directly at:  
**https://redes-proyecto-mcp-673347929287.us-central1.run.app/mcp**

---

## 📦 Installation

1. **Clone the repository**

```bash
git clone https://github.com/Gerax5/Redes_Proyecto1.git
cd RemoteServer
```

2. **Install dependencies**

```bash
pip install fastapi uvicorn requests
```

3. **Start the server**

```bash
uvicorn server:app --reload --port 8080
```

## 🧐 Integration

To use this MCP in your client, add the following block to your configuration file:

```toml
[mcp]
[[mcp.servers]]
name = "pokemon"
type = "http"
url = "http://localhost:8080/mcp"
```

## 🛠️ Endpoints

The server exposes the following endpoints:

### `GET /`

- Returns a simple health check message:
  ```json
  { "message": "Hello from Pokemon MCP!" }
  ```

### `POST /mcp`

Implements the MCP protocol with JSON-RPC methods:

- `initialize`
  Initializes the MCP client-server handshake.

  ````json
  {
  "method": "initialize",
  "id": 1
  }

      ```

  ````

- `tools/list`
  Lists the available tools. In this case, it exposes:
  - get_pokemon: Lookup detailed Pokémon stats and types.
- `tools/call`
  Executes a tool. For example, to retrieve info about Pikachu:

  ```json
  {
    "method": "tools/call",
    "id": 2,
    "params": {
      "arguments": {
        "name": "pikachu"
      }
    }
  }
  ```

  Response example:

  ```json
  {
    "jsonrpc": "2.0",
    "id": 2,
    "result": {
      "content": [
        {
          "type": "text",
          "text": "Pokémon Pikachu (ID 25): Types: electric, Stats: {...}"
        }
      ]
    }
  }
  ```
