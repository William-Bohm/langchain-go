package promptSchema

import "github.com/William-Bohm/langchain-go/langchain-go/rootSchema"

type PromptValue interface {
	ToString() string
	ToMessages() []rootSchema.BaseMessage
}
