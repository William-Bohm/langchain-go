package chains

import (
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/prompt/promptSchema"
)

type ConditionalFunc func(llm llmSchema.BaseLanguageModel) bool

type Conditional struct {
	Condition ConditionalFunc
	Prompt    promptSchema.BasePromptTemplate
}

type BasePromptSelector interface {
	GetPrompt(llm llmSchema.BaseLanguageModel) promptSchema.BasePromptTemplate
}

type ConditionalPromptSelector struct {
	DefaultPrompt promptSchema.BasePromptTemplate
	Conditionals  []Conditional
}

func (cps *ConditionalPromptSelector) GetPrompt(llm llmSchema.BaseLanguageModel) promptSchema.BasePromptTemplate {
	for _, conditional := range cps.Conditionals {
		if conditional.Condition(llm) {
			return conditional.Prompt
		}
	}
	return cps.DefaultPrompt
}

func IsLLM(llm llmSchema.BaseLanguageModel) bool {
	_, ok := llm.(llmSchema.BaseLLM)
	return ok
}

func IsChatModel(llm llmSchema.BaseLanguageModel) bool {
	_, ok := llm.(BaseChatModel)
	return ok
}
