package outputParser

import (
	"github.com/William-Bohm/langchain-go/langchain-go/outputParser/outputParserSchema"
	"strings"
)

type NoOutputParser struct {
	NoOutputStr string
}

func (p *NoOutputParser) Parse(text string) string {
	cleanedText := strings.TrimSpace(text)
	if cleanedText == p.NoOutputStr {
		return ""
	}
	return cleanedText
}

func (p *NoOutputParser) ParseWithPrompt(completion string, prompt outputParserSchema.PromptValue) (interface{}, error) {
	return p.Parse(completion), nil
}

func (p *NoOutputParser) GetFormatInstructions() string {
	return "not implemented"
}

func (p *NoOutputParser) Type() string {
	return "NoOutputParser"
}

func (p *NoOutputParser) ToDict() (map[string]interface{}, error) {
	return map[string]interface{}{
		"_type":         p.Type(),
		"no_output_str": p.NoOutputStr,
	}, nil
}
