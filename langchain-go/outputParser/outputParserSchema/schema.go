package outputParserSchema

import "github.com/William-Bohm/langchain-go/langchain-go/rootSchema"

type BaseOutputParser interface {
	Parse(text string) (interface{}, error)
	ParseWithPrompt(completion string, prompt PromptValue) (interface{}, error)
	GetFormatInstructions() string
	Type() string
	ToDict() (map[string]interface{}, error)
}

type PromptValue interface {
	ToString() string
	ToMessages() []rootSchema.BaseMessageInterface
}
