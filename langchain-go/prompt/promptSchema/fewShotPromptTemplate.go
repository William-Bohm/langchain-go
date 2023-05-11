package promptSchema

import (
	"errors"
	"fmt"
	"strings"

	"github.com/William-Bohm/langchain-go/langchain-go/prompt/exampleSelector" // Replace with the actual package name
)

type FewShotPromptTemplate struct {
	BasePromptTemplate
	Examples         []map[string]interface{}
	ExampleSelector  exampleSelector.BaseExampleSelector
	Suffix           string
	ExampleSeparator string
	Prefix           string
	TemplateFormat   string
	ValidateTemplate bool
}

func NewFewShotPromptTemplate() *FewShotPromptTemplate {
	return &FewShotPromptTemplate{
		ExampleSeparator: "\n\n",
		Prefix:           "",
		TemplateFormat:   "f-string",
		ValidateTemplate: true,
	}
}

func (fspt *FewShotPromptTemplate) CheckExamplesAndSelector() error {
	examples := fspt.Examples != nil
	exampleSelector := fspt.ExampleSelector != nil

	if examples && exampleSelector {
		return errors.New("Only one of 'examples' and 'example_selector' should be provided")
	}

	if !examples && !exampleSelector {
		return errors.New("One of 'examples' and 'example_selector' should be provided")
	}

	return nil
}

func (fspt *FewShotPromptTemplate) TemplateIsValid(partialVariables []string) error {
	if fspt.ValidateTemplate {
		return CheckValidTemplate(
			fspt.Prefix+fspt.Suffix,
			fspt.TemplateFormat,
			append(fspt.InputVariables, partialVariables...),
		)
	}
	return nil
}

func (fspt *FewShotPromptTemplate) GetExamples(kwargs map[string]interface{}) ([]map[string]interface{}, error) {
	if fspt.Examples != nil {
		return fspt.Examples, nil
	} else if fspt.ExampleSelector != nil {
		return fspt.ExampleSelector.SelectExamples(kwargs), nil
	} else {
		return nil, errors.New("No examples or example selector provided")
	}
}

func (fspt *FewShotPromptTemplate) Format(kwargs map[string]interface{}) (string, error) {
	kwargs = fspt.MergePartialAndUserVariables(kwargs)

	examples, err := fspt.GetExamples(kwargs)
	if err != nil {
		return "", err
	}

	exampleStrings := make([]string, len(examples))
	for i, example := range examples {
		formattedExample, err := fspt.Format(example)
		if err != nil {
			return "", err
		}
		exampleStrings[i] = formattedExample
	}

	pieces := append([]string{fspt.Prefix}, exampleStrings...)
	pieces = append(pieces, fspt.Suffix)
	template := strings.Join(pieces, fspt.ExampleSeparator)

	return ExecuteFormatter(template, fspt.TemplateFormat, fspt.InputVariables), nil
}

func (fspt *FewShotPromptTemplate) PromptType() string {
	return "few_shot"
}

func (fspt *FewShotPromptTemplate) ToMap() (map[string]interface{}, error) {
	if fspt.ExampleSelector != nil {
		return nil, errors.New("Saving an example selector is not currently supported")
	}

	return fspt.ToMap(), nil
}

func main() {
	// Usage example
	template := NewFewShotPromptTemplate()
	err := template.CheckExamplesAndSelector()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = template.TemplateIsValid()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	formatted, err := template.Format(map[string]interface{}{"variable1": "foo"})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Formatted template:", formatted)

	// Example usage of ToMap
	m, err := template.ToMap()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Template as map:", m)
}
