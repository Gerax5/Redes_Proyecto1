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

	ant "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/anthropics/anthropic-sdk-go/packages/param"

)

func ConnectMCPServer(server config.MCPServer) (*client.Client, error) {
	stdio := transport.NewStdio(server.Command, os.Environ(), server.Args...)
	c := client.NewClient(stdio)

	ctx := context.Background()

	if err := c.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start transport: %w", err)
	}

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
	cfg, errConfig := config.Load("config.toml")
	if errConfig != nil {
		log.Fatal("Error cargando config.toml:", errConfig)
		return
	}

	fmt.Printf("Config: %+v\n", cfg)


	var mcpServers []*client.Client
	var toolsDesc []string
	var toolsForClaude []ant.ToolUnionParam 
	ctx := context.Background()

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
				}
		}
	}

	// debug config
	b, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(b))

	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar .env, usando variables de entorno del sistema")
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("Falta ANTHROPIC_API_KEY en el entorno")
	}

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

		history = append(history, ant.NewUserMessage(ant.NewTextBlock(input)))

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

		for _, block := range resp.Content {
			if t := block.AsText(); t.Text != "" {
				// respuesta normal de texto
				fmt.Println("Claude:", t.Text)
				history = append(history, ant.NewAssistantMessage(ant.NewTextBlock(t.Text)))
			}

			if toolUse := block.AsToolUse(); toolUse.Type == "tool_use" {
				fmt.Printf("Claude quiere usar tool: %s con input: %+v\n", toolUse.Name, toolUse.Input)

				for _, c := range mcpServers {
					log.Println("➡️  Llamando a CallTool en servidor MCP...")
					// callRes, err := c.CallTool(ctx, mcp.CallToolRequest{
					// 	Params: mcp.CallToolParams{
					// 		Name:      toolUse.Name,
					// 		Arguments: toolUse.Input,
					// 	},
					// })
					// if err != nil {
					// 	log.Printf("error ejecutando tool %s: %v", toolUse.Name, err)
					// 	continue
					// }
					req := mcp.CallToolRequest{
    Params: mcp.CallToolParams{
        Name:      toolUse.Name,
        Arguments: toolUse.Input,
    },
}
reqJSON, _ := json.MarshalIndent(req, "", "  ")
log.Printf("➡️  Enviando CallTool a:\n%s", string(reqJSON))

callRes, err := c.CallTool(ctx, req)
if err != nil {
    log.Printf("❌ Error en CallTool: %v", err)
    continue
}

					toolResultText := fmt.Sprintf("Tool %s result: %+v", toolUse.Name, callRes)
					history = append(history, ant.NewAssistantMessage(ant.NewTextBlock(toolResultText)))

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
						}
					}
				}
			}
		}

	}
}
