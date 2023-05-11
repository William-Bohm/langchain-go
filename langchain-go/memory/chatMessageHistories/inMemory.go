package chatMessageHistories

import (
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
)

type ChatMessageHistory struct {
	messages []rootSchema.BaseMessageInterface
}

func NewChatMessageHistory() *ChatMessageHistory {
	return &ChatMessageHistory{
		messages: []rootSchema.BaseMessageInterface{},
	}
}

func (c *ChatMessageHistory) Messages() ([]rootSchema.BaseMessageInterface, error) {
	return c.messages, nil
}

func (c *ChatMessageHistory) AddUserMessage(message string) error {
	humanMessage := rootSchema.NewHumanMessage(message)
	c.messages = append(c.messages, humanMessage)
	return nil
}

func (c *ChatMessageHistory) AddAIMessage(message string) error {
	aiMessage := rootSchema.NewAIMessage(message)
	c.messages = append(c.messages, aiMessage)
	return nil
}

func (c *ChatMessageHistory) Clear() error {
	c.messages = []rootSchema.BaseMessageInterface{}
	return nil
}
