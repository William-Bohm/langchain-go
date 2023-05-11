package outputParser

import (
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/outputParser/outputParserSchema"
	"strings"
)

// BooleanOutputParser represents the boolean output parser
type BooleanOutputParser struct {
	TrueVal  string
	FalseVal string
}

// Parse parses the output of an LLM call to a boolean
func (p BooleanOutputParser) Parse(text string) (bool, error) {
	cleanedText := strings.TrimSpace(text)
	if cleanedText != p.TrueVal && cleanedText != p.FalseVal {
		return false, errors.New("BooleanOutputParser expected output value to either be " + p.TrueVal + " or " + p.FalseVal + ". Received " + cleanedText)
	}
	return cleanedText == p.TrueVal, nil
}

// ParseWithPrompt parses the output of an LLM call with a prompt
func (p BooleanOutputParser) ParseWithPrompt(completion string, prompt outputParserSchema.PromptValue) (bool, error) {
	return p.Parse(completion)
}

// GetType returns the snake-case string identifier for output parser type
func (p BooleanOutputParser) GetType() string {
	return "boolean_output_parser"
}

// GetFormatInstructions returns instructions on how the LLM output should be formatted
func (p BooleanOutputParser) GetFormatInstructions() string {
	return "Not Implemented"
}

// Dict returns a map representation of the output parser
func (p BooleanOutputParser) Dict() (map[string]interface{}, error) {
	outputParserDict := make(map[string]interface{})
	outputParserDict["TrueVal"] = p.TrueVal
	outputParserDict["FalseVal"] = p.FalseVal
	outputParserDict["_type"] = p.GetType()
	return outputParserDict, nil
}
