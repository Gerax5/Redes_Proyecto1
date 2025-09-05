package main

import (
    "bufio"
    "context"
    "fmt"
    "log"
    "os"
	"encoding/json"
    "strings"
	"github.com/joho/godotenv"
	"proyecto/config"

    ant "github.com/anthropics/anthropic-sdk-go"
    "github.com/anthropics/anthropic-sdk-go/option"
)


func main() {
	cfg, _ := config.Load("config.toml")

	b, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(b))

	if err := godotenv.Load(); err != nil {
        log.Println("No se pudo cargar .env, usando variables de entorno del sistema")
    }

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ Falta ANTHROPIC_API_KEY en el entorno")
	}

	claude := ant.NewClient(option.WithAPIKey(apiKey))
	ctx := context.Background()
	
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
		})
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		var answer string
		for _, c := range resp.Content {
			if t := c.AsText(); t.Text != "" {
				answer = t.Text
			}
		}

		fmt.Println("Claude:", answer)

		history = append(history, ant.NewAssistantMessage(ant.NewTextBlock(answer)))
	}
}
