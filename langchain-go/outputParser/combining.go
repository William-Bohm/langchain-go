package outputParser

import (
	"errors"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/outputParser/outputParserSchema"
	"strings"
)

// CombiningOutputParser combines multiple output parsers into one
type CombiningOutputParser struct {
	Parsers []outputParserSchema.BaseOutputParser
}

// NewCombiningOutputParser creates a new CombiningOutputParser
func NewCombiningOutputParser(parsers []outputParserSchema.BaseOutputParser) (*CombiningOutputParser, error) {
	if len(parsers) < 2 {
		return nil, errors.New("Must have at least two parsers")
	}
	for _, parser := range parsers {
		if parser.Type() == "combining" {
			return nil, errors.New("Cannot nest combining parsers")
		}
		if parser.Type() == "list" {
			return nil, errors.New("Cannot combine list parsers")
		}
	}
	return &CombiningOutputParser{Parsers: parsers}, nil
}

// Parse parses the output of an LLM call
func (p CombiningOutputParser) Parse(text string) (map[string]interface{}, error) {
	texts := strings.Split(text, "\n\n")
	output := make(map[string]interface{})
	for i, parser := range p.Parsers {
		result, err := parser.Parse(strings.TrimSpace(texts[i]))
		if err != nil {
			return nil, err
		}
		// Assuming that result is of type map[string]interface{}
		// you might need to do type assertion depending on your actual implementation
		for k, v := range result.(map[string]interface{}) {
			output[k] = v
		}
	}
	return output, nil
}

// ParseWithPrompt parses the output of an LLM call with a prompt
func (p CombiningOutputParser) ParseWithPrompt(completion string, prompt outputParserSchema.PromptValue) (interface{}, error) {
	return p.Parse(completion)
}

// Type returns the snake-case string identifier for output parser type
func (p CombiningOutputParser) Type() string {
	return "combining"
}

// GetFormatInstructions returns instructions on how the LLM output should be formatted
func (p CombiningOutputParser) GetFormatInstructions() string {
	initial := fmt.Sprintf("For your first output: %s", p.Parsers[0].GetFormatInstructions())
	var subsequent []string
	for _, parser := range p.Parsers[1:] {
		subsequent = append(subsequent, fmt.Sprintf("Complete that output fully. Then produce another output, separated by two newline characters: %s", parser.GetFormatInstructions()))
	}
	return fmt.Sprintf("%s\n%s", initial, strings.Join(subsequent, "\n"))
}

// ToDict returns a map representation of the output parser
func (p CombiningOutputParser) ToDict() (map[string]interface{}, error) {
	outputParserDict := make(map[string]interface{})
	var parserTypes []string
	for _, parser := range p.Parsers {
		parserTypes = append(parserTypes, parser.Type())
	}
	outputParserDict["Parsers"] = parserTypes
	outputParserDict["_type"] = p.Type()
	return outputParserDict, nil
}
