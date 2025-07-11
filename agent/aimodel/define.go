package aimodel

import (
	"context"
	"fmt"
	"strings"
)

type UserContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

const DoneKey = "[DONE]"

type EAIChatRole string

const (
	EAIChatRoleUser      EAIChatRole = "user"
	EAIChatRoleSystem    EAIChatRole = "system"
	EAIChatRoleAssistant EAIChatRole = "assistant"
)

type Chat struct {
	Role    EAIChatRole `json:"role"`
	Content any         `json:"content"`
}

type ChatList []Chat

func (c ChatList) ToString() string {
	str := strings.Builder{}
	for _, chat := range c {
		str.WriteString(fmt.Sprintf("%s:\n %s\n", chat.Role, chat.Content))
	}
	return str.String()
}

type IAiModel interface {
	Chat(ctx context.Context, chatList ChatList) (StringResult, error)
	ChatFlow(ctx context.Context, chatData ChatList) (resChan chan string, thinkChan chan string, errChan chan error, err error)
}

type Result[T any] struct {
	ID     string `json:"id"`
	Result T      `json:"result"`
	Reason string `json:"reason,omitempty"`
}

type StringResult = Result[string]

type OpenAIResponse struct {
	Choices []OpenAIChoices `json:"choices"`
	ID      string          `json:"id"`
}

type OpenAIChoices struct {
	FinishReason interface{}        `json:"finish_reason"`
	Index        int                `json:"index"`
	Delta        OpenAIChoicesDelta `json:"delta"`
}

type OpenAIChoicesDelta struct {
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
}
