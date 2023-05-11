package chatMessageHistories

import (
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
)

type BaseChatMessageHistory interface {
	Messages() ([]rootSchema.BaseMessageInterface, error)
	AddUserMessage(message string) error
	AddAIMessage(message string) error
	Clear() error
}
