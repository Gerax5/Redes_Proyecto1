package llm

import (
    "bytes"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "os"
)

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type AnthropicRequest struct {
    Model     string    `json:"model"`
    MaxTokens int       `json:"max_tokens"`
    Messages  []Message `json:"messages"`
}

type AnthropicResponse struct {
    Content []struct {
        Text string `json:"text"`
    } `json:"content"`
}

func QueryAnthropic(history []Message) (string, error) {
    url := "https://api.anthropic.com/v1/messages"

    body := AnthropicRequest{
        Model:     "claude-3-haiku-20240307",
        MaxTokens: 300,
        Messages:  history,
    }

    jsonBody, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    apiKey := os.Getenv("ANTHROPIC_API_KEY")
    req.Header.Set("x-api-key", apiKey)
    req.Header.Set("anthropic-version", "2023-06-01")


    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    responseBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var response AnthropicResponse
    if err := json.Unmarshal(responseBytes, &response); err != nil {
        return "", err
    }

    if len(response.Content) > 0 {
        return response.Content[0].Text, nil
    }

    return "(sin respuesta)", nil
}
