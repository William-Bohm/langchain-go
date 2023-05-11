package llmSchema

import (
	"github.com/William-Bohm/langchain-go/langchain-go/llm/openai"
)

var LLMTypeToClassMap = make(map[string]func(map[string]interface{}) (BaseLanguageModel, error))

func init() {
	LLMTypeToClassMap["openai"] = openai.New
}
