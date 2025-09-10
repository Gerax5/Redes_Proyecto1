# 🤖 MCP + Claude Integration

This project provides a command-line interface that connects to **MCP servers** (both HTTP and stdio) and integrates them with **Anthropic's Claude** model.  
The program allows interactive conversations with Claude, including the ability for Claude to request the execution of tools provided by MCP servers.

---

## 🪶 Features

- Load and validate configuration from a `config.toml` file.
- Connect to multiple MCP servers (HTTP and stdio supported).
- Retrieve and expose tools from MCP servers to Claude.
- Interactive CLI interface where the user can chat with Claude.
- Claude can request tool usage, which is executed automatically.
- Chat history is stored in a `chat.log` file.
- Environment variables are loaded from `.env` (fallback to system env).

---

## 🤨 Requirements

- **Go 1.21+**
- [Anthropic API Key](https://docs.anthropic.com/) (`ANTHROPIC_API_KEY`)
- Dependencies (managed via `go.mod`):
  - `github.com/joho/godotenv`
  - `github.com/BurntSushi/toml`
  - `github.com/anthropics/anthropic-sdk-go`
  - `github.com/mark3labs/mcp-go`

---

## 🙏 Install Go 🙏

Before running the project, you must install **Go** (version 1.21 or higher):

1. **Download Go from the official website:**  
   [https://go.dev/dl/](https://go.dev/dl/)

2. **Follow the installation instructions for your operating system:**

   - **Linux/macOS**: Extract the archive to `/usr/local` and add Go to your PATH:
     ```bash
     tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
     export PATH=$PATH:/usr/local/go/bin
     ```
   - **Windows**: Run the `.msi` installer and follow the wizard.

3. **Verify the installation:**
   ```bash
   go version
   ```
   You should see something like:
   ```bash
   go version go1.21.0 linux/amd64
   ```

---

## 📦 Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Gerax5/MCP_Server
   cd MCP_Server
   ```

2. **Initialize and update submodules:**
   ````bash`
   git submodule update --init --recursive

   ```

   ```

3. **Install dependencies:**

   ```bash
   go mod tidy
   ```

4. **Create a `.env` file with your Anthropic API key:**

   ```bash
   ANTHROPIC_API_KEY=your_api_key_here
   ```

5. **Create a `config.toml` file to define MCP servers:**

   ```bash
   [mcp]
   [[mcp.servers]]
   name = "youtube"
   type = "stdio"
   command = "Python3"
   args  = ["./MCP_Server/Server.py"]

   [[mcp.servers]]
   name = "pokemon"
   type = "http"
   url = "https://redes-proyecto-mcp-673347929287.us-central1.run.app/mcp"
   ```

---

## 🌵 Usage

Run the program:

```bash
go run main.go
```

---

## 📩 Logs

All interactions are saved in a file called `chat.log` with timestamps:

```bash
[2025-09-08T12:34:56Z]
User: Hello
Bot: Hi there! How can I help you today?
```
