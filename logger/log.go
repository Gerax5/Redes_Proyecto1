// Package logger provides a simple utility for saving chat logs.
package logger

import (
    "fmt"
    "os"
    "time"
)
// SaveLog appends a new log entry containing user input and bot response
// into a file named chat.log with a timestamp.
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