# 🐁 Pokémon (hol demonio) MCP

Este proyecto implementa un servidor MCP que expone la PokeAPI como herramienta remota para integrarse con un chatbot.

## ☁️ Usa el deploy

Server, el mcp esta en el endpint `/mcp`  
**https://redes-proyecto-mcp-673347929287.us-central1.run.app**

usar:  
**https://redes-proyecto-mcp-673347929287.us-central1.run.app/mcp**

## 📦 Instalación

1. _Clonar Repositorio_

```bash
git clone https://github.com/Gerax5/Redes_Proyecto1.git
cd RemoteServer
```

2. _Instalar depedencias_

```bash
pip install fastapi uvicorn requests
```

3. _Inicia el servidor_

```bash
uvicorn server:app --reload --port 8080
```

## 🧐 Implementarlo

Para usar este MCP en tu cliente:

```json
[mcp]
[[mcp.servers]]
name = "pokemon"
type = "http"
url = "http://localhost:8080/mcp"
```
