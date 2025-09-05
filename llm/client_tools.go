package llm

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
)


type ContentBlock struct {
    Type  string                 `json:"type"`
    Text  string                 `json:"text,omitempty"`
    ID    string                 `json:"id,omitempty"`
    Name  string                 `json:"name,omitempty"`
    Input map[string]interface{} `json:"input,omitempty"`
}

type ToolMessage struct {
    Role    string         `json:"role"`
    Content []ContentBlock `json:"content"`
}

type ToolRequest struct {
    Model     string        `json:"model"`
    MaxTokens int           `json:"max_tokens"`
    Messages  []ToolMessage `json:"messages"`
    Tools     []any         `json:"tools,omitempty"`
}

type ToolResponse struct {
    Content []ContentBlock `json:"content"`
}


func QueryAnthropicWithTools(apiKey string, history []ToolMessage, tools []any) (ToolResponse, error) {
    url := "https://api.anthropic.com/v1/messages"

    body := ToolRequest{
        Model:     "claude-3-haiku-20240307",
        MaxTokens: 300,
        Messages:  history,
        Tools:     tools,
    }

    jsonBody, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-api-key", apiKey)
    req.Header.Set("anthropic-version", "2023-06-01")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return ToolResponse{}, err
    }
    defer resp.Body.Close()

    responseBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return ToolResponse{}, err
    }

    var response ToolResponse
    if err := json.Unmarshal(responseBytes, &response); err != nil {
        return ToolResponse{}, err
    }

    return response, nil
}
