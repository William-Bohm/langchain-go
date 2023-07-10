package outputParser

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/chains"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/outputParser/outputParserSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/prompt/promptSchema"
)

// OutputFixingParser represents a parser that tries to fix parsing errors
type OutputFixingParser struct {
	Parser     outputParserSchema.BaseOutputParser
	RetryChain chains.LLMChain // Assuming you have this LLMChain defined elsewhere
}

// FromLLM creates an instance of OutputFixingParser from an LLM
func FromLLM(llm llmSchema.BaseLanguageModel, parser outputParserSchema.BaseOutputParser, prompt promptSchema.BasePromptTemplate) OutputFixingParser {
	chain := chains.NewLLMChain(llm, prompt) // Assuming you have NewLLMChain function defined elsewhere
	return OutputFixingParser{
		Parser:     parser,
		RetryChain: chain,
	}
}

// Parse tries to parse the completion, and retries in case of an error
func (p OutputFixingParser) Parse(completion string) (interface{}, error) {
	parsedCompletion, err := p.Parser.Parse(completion)
	if err != nil {
		newCompletion, _ := p.RetryChain.Run(
			p.Parser.GetFormatInstructions(),
			completion,
			fmt.Sprintf("%v", err),
		)
		parsedCompletion, err = p.Parser.Parse(newCompletion)
	}
	return parsedCompletion, err
}

// ParseWithPrompt tries to parse the completion with a prompt, and retries in case of an error
func (p OutputFixingParser) ParseWithPrompt(completion string, prompt outputParserSchema.PromptValue) (interface{}, error) {
	return p.Parse(completion)
}

// GetFormatInstructions returns the format instructions of the internal parser
func (p OutputFixingParser) GetFormatInstructions() string {
	return p.Parser.GetFormatInstructions()
}

// GetType returns the type of the internal parser
func (p OutputFixingParser) GetType() string {
	return p.Parser.Type()
}

// Dict returns a map representation of the output parser
func (p OutputFixingParser) Dict() (map[string]interface{}, error) {
	outputParserDict := make(map[string]interface{})
	outputParserDict["_type"] = p.GetType()
	// Add more fields as necessary
	return outputParserDict, nil
}
