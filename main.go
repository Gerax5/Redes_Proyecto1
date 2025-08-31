package main

import (
    "bufio"
    "fmt"
    "os"
    "log"
    "github.com/joho/godotenv"
    "strings"
    "proyecto/llm"
    "proyecto/session"
    "proyecto/logger"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No se pudo cargar .env, usando variables de entorno del sistema")
    }
    fmt.Println("Chatbot MCP!")
    fmt.Println("Escribe 'exit' para salir.")
    fmt.Println("Comando especial: youtube <channel_id> para analizar comunidad")

    s := session.NewSession()
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print(">>> ")
        scanner.Scan()
        input := strings.TrimSpace(scanner.Text())

        if input == "exit" {
            break
        }

        s.AddMessage("user", input)

        response, err := llm.QueryAnthropic(s.GetHistory())
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

        s.AddMessage("assistant", response)

        fmt.Println("Claude:", response)

        logger.SaveLog(input, response)
    }
}
