package session

import "proyecto/llm"

type Session struct {
    history []llm.Message
}

func NewSession() *Session {
    return &Session{history: []llm.Message{}}
}

func (s *Session) AddMessage(role, content string) {
    s.history = append(s.history, llm.Message{Role: role, Content: content})
}

func (s *Session) GetHistory() []llm.Message {
    return s.history
}
