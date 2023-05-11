package rootSchema

import (
	"github.com/William-Bohm/langchain-go/langchain-go/prompt/promptSchema"
	"strings"
)

type BaseLanguageModelInterface interface {
	GeneratePrompt(prompts []promptSchema.PromptValue, stop []string) (LLMResult, error)
	AGeneratePrompt(prompts []promptSchema.PromptValue, stop []string) (LLMResult, error)
	GetNumTokens(text string) int
	GetNumTokensFromMessages(messages []BaseMessageInterface) int
}

type BaseLanguageModel struct {
}

func (b *BaseLanguageModel) GetNumTokens(text string) int {
	// This function is a simple implementation assuming each word is a token.
	// You may need to replace it with a library that can accurately tokenize text similar to the GPT-3 tokenizer
	tokens := strings.Fields(text)
	return len(tokens)
}

func (b *BaseLanguageModel) GetNumTokensFromMessages(messages []BaseMessageInterface) int {
	// Implementation goes here
}

type Generation struct {
	Text           string                 `json:"text"`                      // Generated text output
	GenerationInfo map[string]interface{} `json:"generation_info,omitempty"` // Raw generation info response from the provider
}

type LLMResult struct {
	Generations []Generation           `json:"generations"`          // List of the things generated
	LLMOutput   map[string]interface{} `json:"llm_output,omitempty"` // For arbitrary LLM provider specific output
}
