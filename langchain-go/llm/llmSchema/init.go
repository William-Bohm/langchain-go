package llmSchema

import (
	"github.com/William-Bohm/langchain-go/langchain-go/llm/openai"
)

// Declare a type that matches your desired function signature
type LLMFactoryFunc func(...openai.Option) (BaseLanguageModel, error)

var LLMTypeToClassMap = make(map[string]LLMFactoryFunc)

func init() {
	// Use a type conversion to make `openai.New` match the desired function signature
	LLMTypeToClassMap["openai"] = LLMFactoryFunc(func(options ...openai.Option) (BaseLanguageModel, error) {
		llm, err := openai.New(options...)
		// Converting *OpenaiLLM to *BaseLanguageModel
		return BaseLanguageModel(llm), err
	})
}
