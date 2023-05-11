package rootSchema

import (
	"fmt"
	"strings"
)

/*
*
*
*
* Base Messages
*
*
 */
type MessageType string

const (
	Human  MessageType = "human"
	AI     MessageType = "ai"
	System MessageType = "system"
	Chat   MessageType = "chat"
	Base   MessageType = "base"
)

type BaseMessageInterface interface {
	Type() MessageType
}

type BaseMessage struct {
	Content          string
	AdditionalKwargs map[string]interface{}
}

func NewBaseMessage(content string, messageType MessageType) *BaseMessage {
	return &BaseMessage{
		Content:          content,
		AdditionalKwargs: make(map[string]interface{}),
	}
}

func (bm *BaseMessage) GetContent() string {
	return bm.Content
}

func (bm *BaseMessage) GetAdditionalKwargs() map[string]interface{} {
	return bm.AdditionalKwargs
}

type HumanMessage struct {
	BaseMessage
}

func NewHumanMessage(content string) *HumanMessage {
	return &HumanMessage{
		BaseMessage: *NewBaseMessage(content, Human),
	}
}

func (hm *HumanMessage) Type() MessageType {
	return Human
}

type AIMessage struct {
	BaseMessage
}

func NewAIMessage(content string) *AIMessage {
	return &AIMessage{
		BaseMessage: *NewBaseMessage(content, AI),
	}
}

func (ai *AIMessage) Type() MessageType {
	return AI
}

type SystemMessage struct {
	BaseMessage
}

func NewSystemMessage(content string) *SystemMessage {
	return &SystemMessage{
		BaseMessage: *NewBaseMessage(content, System),
	}
}

func (sm *SystemMessage) Type() MessageType {
	return System
}

type ChatMessage struct {
	BaseMessage
	Role string
}

func NewChatMessage(content string, role string) *ChatMessage {
	return &ChatMessage{
		BaseMessage: *NewBaseMessage(content, Chat),
		Role:        role,
	}
}

func (cm *ChatMessage) Type() MessageType {
	return Chat
}

func GetBufferString(messages []BaseMessageInterface, prefixes ...string) (string, error) {
	humanPrefix := "Human"
	aiPrefix := "AI"

	if len(prefixes) > 0 {
		humanPrefix = prefixes[0]
	}
	if len(prefixes) > 1 {
		aiPrefix = prefixes[1]
	}

	var stringMessages []string

	for _, m := range messages {
		var role string
		switch msg := m.(type) {
		case *HumanMessage:
			role = humanPrefix
		case *AIMessage:
			role = aiPrefix
		case *SystemMessage:
			role = "System"
		case *ChatMessage:
			role = msg.Role
		default:
			return "", fmt.Errorf("got unsupported message type: %v", m)
		}
		stringMessages = append(stringMessages, fmt.Sprintf("%s: %s", role, m.Content()))
	}

	return strings.Join(stringMessages, "\n"), nil
}
