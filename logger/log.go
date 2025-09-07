package logger

import (
    "fmt"
    "os"
    "time"
)

func SaveLog(input, response string) {
    f, err := os.OpenFile("chat.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error abriendo log:", err)
        return
    }
    defer f.Close()

    log := fmt.Sprintf("[%s]\nUser: %s\nBot: %s\n\n",
        time.Now().Format(time.RFC3339), input, response)

    f.WriteString(log)
}