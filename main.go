// This program connects to multiple MCP servers (both HTTP and stdio),
// integrates them with Anthropic's Claude model, and provides a CLI
// where the user can interact with Claude. If Claude requests the use
// of a tool, the program executes it through MCP servers and returns
// the result back to Claude, maintaining a conversation history.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"proyecto/config"
	"proyecto/logger"

	ant "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/anthropics/anthropic-sdk-go/packages/param"

)

// ConnectMCPServer establishes a connection to a stdio MCP server.
func ConnectMCPServer(server config.MCPServer) (*client.Client, error) {
	stdio := transport.NewStdio(server.Command, os.Environ(), server.Args...)
	c := client.NewClient(stdio)

	ctx := context.Background()

	if err := c.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start transport: %w", err)
	}

	// Initialize MCP session
	initReq := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "my-mcp-client",
				Version: "0.1.0",
			},
			Capabilities: mcp.ClientCapabilities{},
		},
	}

	_, err := c.Initialize(ctx, initReq)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize: %w", err)
	}

	return c, nil
}

// ConnectMCPHTTPServer establishes a connection to an HTTP MCP server.
func ConnectMCPHTTPServer(server config.MCPServer) (*client.Client, error) {
    trans, err := transport.NewStreamableHTTP(
				server.URL,
			)
	if err != nil {
		fmt.Printf("Failed to connect to `%s`: %s", server.URL, err)
	}
    c := client.NewClient(trans)

    ctx := context.Background()
    if err := c.Start(ctx); err != nil {
        return nil, fmt.Errorf("failed to start HTTP transport: %w", err)
    }

	// Initialize MCP session
    initReq := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "my-mcp-client",
				Version: "0.1.0",
			},
			Capabilities: mcp.ClientCapabilities{},
		},
	}

    _, err2 := c.Initialize(ctx, initReq)
    if err2 != nil {
        return nil, fmt.Errorf("failed to initialize HTTP MCP server: %w", err2)
    }

    return c, nil
}


func main() {
	// Load configuration from config.toml
	cfg, errConfig := config.Load("config.toml")
	if errConfig != nil {
		log.Fatal("Error cargando config.toml:", errConfig)
		return
	}

	var mcpServers []*client.Client
	var toolsDesc []string
	var toolsForClaude []ant.ToolUnionParam 
	toolToServer := map[string]*client.Client{}
	ctx := context.Background()

	// Connect to each MCP server and retrieve available tools
	for _, server := range cfg.MCP.Servers {
		fmt.Printf("Servidor MCP: %s (%s)\n", server.Name, server.Type)

		switch server.Type {
		case "http":
			fmt.Printf(" -> URL: %s\n", server.URL)
				client, err := ConnectMCPHTTPServer(server)
				if err != nil {
					log.Fatal(err)
				}
				mcpServers = append(mcpServers, client)

				// List available tools
				toolsRes, err := client.ListTools(ctx, mcp.ListToolsRequest{})
				if err != nil {
					log.Printf("error enviando tools/list a %s: %v", server.Name, err)
					continue
				}
				for _, t := range toolsRes.Tools {
					toolsDesc = append(toolsDesc,
						fmt.Sprintf("- %s: %s", t.Name, t.Description))

					toolsForClaude = append(toolsForClaude, ant.ToolUnionParam{
						OfTool: &ant.ToolParam{
							Name:        t.Name,
							Description: param.NewOpt(t.Description),
							InputSchema: ant.ToolInputSchemaParam{
								Required:   t.InputSchema.Required,
								Properties: t.InputSchema.Properties,
							},
						},
					})
					
					toolToServer[t.Name] = client

				}

		case "stdio":
			client, err := ConnectMCPServer(server)
			if err != nil {
				log.Fatal(err)
			}
			mcpServers = append(mcpServers, client)

			toolsRes, err := client.ListTools(ctx, mcp.ListToolsRequest{})
			if err != nil {
				log.Printf("error enviando tools/list a %s: %v", server.Name, err)
				continue
			}
			for _, t := range toolsRes.Tools {
				toolsDesc = append(toolsDesc,
					fmt.Sprintf("- %s: %s", t.Name, t.Description))

				// List available tools
				toolsForClaude = append(toolsForClaude, ant.ToolUnionParam{
						OfTool: &ant.ToolParam{
							Name:        t.Name,
							Description: param.NewOpt(t.Description),
							InputSchema: ant.ToolInputSchemaParam{
								Required:   t.InputSchema.Required,
								Properties: t.InputSchema.Properties,
							},
						},	
				})

				toolToServer[t.Name] = client
			}
		}
	}


	// debug config
	b, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(b))

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar .env, usando variables de entorno del sistema")
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("Falta ANTHROPIC_API_KEY en el entorno")
	}
	
	// Initialize Claude client
	claude := ant.NewClient(option.WithAPIKey(apiKey))


	history := []ant.MessageParam{}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">>> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}

		logger.SaveLog(input, "(esperando respuesta de Claude...)")
		history = append(history, ant.NewUserMessage(ant.NewTextBlock(input)))

		// Send request to Claude
		resp, err := claude.Messages.New(ctx, ant.MessageNewParams{
			Model:     ant.Model("claude-3-haiku-20240307"),
			MaxTokens: 300,
			Messages:  history,
			Tools:     toolsForClaude,
		})
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		var answer string
		// Handle Claude response and tool usage
		for _, block := range resp.Content {
			if t := block.AsText(); t.Text != "" {
				answer = t.Text
				fmt.Println("Claude:", answer)
				history = append(history, ant.NewAssistantMessage(ant.NewTextBlock(answer)))
				logger.SaveLog(input, answer)
			}

			// If Claude requests a tool
			if toolUse := block.AsToolUse(); toolUse.Type == "tool_use" {
				inputJSON := string(toolUse.Input)
				fmt.Printf("Claude quiere usar tool: %s con input: %s\n", toolUse.Name, inputJSON)

				for _, c := range mcpServers {
					if c != toolToServer[toolUse.Name] {
						continue
					}

					logger.SaveLog(fmt.Sprintf("CALL tool %s", toolUse.Name), inputJSON)
					callRes, err := c.CallTool(ctx, mcp.CallToolRequest{
						Params: mcp.CallToolParams{
							Name:      toolUse.Name,
							Arguments: toolUse.Input,
						},
					})
					if err != nil {
						log.Printf("error ejecutando tool %s: %v", toolUse.Name, err)
						continue
					}

					toolResultText := fmt.Sprintf("Tool %s result: %+v", toolUse.Name, callRes)
					history = append(history, ant.NewAssistantMessage(ant.NewTextBlock(toolResultText)))
					logger.SaveLog(fmt.Sprintf("RESULT tool %s", toolUse.Name), toolResultText)

					followUp, err := claude.Messages.New(ctx, ant.MessageNewParams{
						Model:     ant.Model("claude-3-haiku-20240307"),
						MaxTokens: 300,
						Messages:  history,
					})
					if err != nil {
						log.Println("Error pidiendo análisis a Claude:", err)
						continue
					}

					for _, c := range followUp.Content {
						if t := c.AsText(); t.Text != "" {
							fmt.Println("Claude (análisis):", t.Text)
							history = append(history, ant.NewAssistantMessage(ant.NewTextBlock(t.Text)))
							logger.SaveLog(fmt.Sprintf("ANALYSIS tool %s", toolUse.Name), t.Text)
						}
					}
				}
			}
		}

	}
}
