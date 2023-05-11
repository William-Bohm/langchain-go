package promptSchema

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
	"strings"
)

type StringPromptValue struct {
	Text string
}

func NewStringPromptValue(text string) *StringPromptValue {
	return &StringPromptValue{Text: text}
}

func (spv *StringPromptValue) ToString() string {
	return spv.Text
}

func (spv *StringPromptValue) ToMessages() []rootSchema.BaseMessageInterface {
	return []rootSchema.BaseMessageInterface{rootSchema.NewHumanMessage(spv.Text)}
}

type StringPromptTemplate struct {
	BasePromptTemplate
}

func (spt *StringPromptTemplate) format(kwargs map[string]interface{}) (string, error) {
	// TODO: Implement actual desired formatting
	prompt := spt.PromptType
	for k, v := range kwargs {
		placeholder := fmt.Sprintf("{%s}", k)
		value := fmt.Sprint(v)
		prompt = strings.ReplaceAll(prompt, placeholder, value)
	}
	return prompt, nil
}

func (spt *StringPromptTemplate) FormatPrompt(kwargs map[string]interface{}) (PromptValue, error) {
	text, err := spt.format(kwargs)
	if err != nil {
		return nil, err
	}
	return StringPromptValue{Text: text}, nil
}

func NewStringPromptTemplate(inputVars []string, outputParser BaseOutputParser, partialVars map[string]interface{}, promptType string) StringPromptTemplate {
	baseTemplate := NewBasePromptTemplate(inputVars, outputParser, partialVars, promptType)
	return StringPromptTemplate{BasePromptTemplate: baseTemplate}
}
